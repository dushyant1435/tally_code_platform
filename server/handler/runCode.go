package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"y/models"
)

const executionTimeout = 5 * time.Second

type execResult struct {
	Stdout   string
	Stderr   string
	Runtime  float64
	TimedOut bool
	ExitCode int
}

// runUserCode prepares a workspace for the given language, optionally compiles,
// and executes the program with the supplied stdin.
// The returned execResult is empty if compile failed; compileLog explains why.
func runUserCode(ctx context.Context, language, code, stdin string) (
	res execResult, compileLog string, compileFail bool, err error,
) {
	p, err := prepareExecution(ctx, language, code)
	if err != nil {
		return execResult{}, "", false, err
	}
	defer p.Cleanup()

	if p.CompileFail {
		return execResult{}, p.CompileLog, true, nil
	}

	cmd := exec.CommandContext(ctx, p.RunArgs[0], p.RunArgs[1:]...)
	cmd.Dir = p.Dir
	cmd.Stdin = strings.NewReader(stdin)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	runErr := cmd.Run()
	elapsed := time.Since(start).Seconds()

	res = execResult{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Runtime: elapsed,
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		res.TimedOut = true
		return res, "", false, nil
	}
	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			res.ExitCode = exitErr.ExitCode()
			return res, "", false, nil
		}
		return res, "", false, runErr
	}
	return res, "", false, nil
}

// RunCode executes the submitted code against every test case for a problem.
// Records a submissions row with the resulting verdict.
func RunCode(w http.ResponseWriter, r *http.Request) {
	c := ClaimsFromContext(r.Context())
	if c == nil {
		httpError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req models.CodeData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	language, _ := Language(req.Language)
	req.Language = language

	db := DB()
	rows, err := db.Query(`SELECT input, output FROM testcases WHERE id=$1`, req.ID)
	if err != nil {
		log.Printf("RunCode query: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch test cases")
		return
	}
	defer rows.Close()

	var (
		status       = models.StatusAccepted
		message      string
		failedAt     int
		totalRuntime float64
		rowCount     int
	)

	for rows.Next() {
		rowCount++
		var input, expected string
		if err := rows.Scan(&input, &expected); err != nil {
			log.Printf("RunCode scan: %v", err)
			httpError(w, http.StatusInternalServerError, "failed to scan test case")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), executionTimeout)
		res, compileLog, compileFail, runErr := runUserCode(ctx, language, req.Code, input)
		cancel()

		if compileFail {
			status, message, failedAt = models.StatusCompileError, truncate(compileLog, 400), 0
			break
		}

		if res.Runtime > totalRuntime {
			totalRuntime = res.Runtime
		}

		switch {
		case runErr != nil:
			status, message, failedAt = models.StatusServerError, runErr.Error(), rowCount
		case res.TimedOut:
			status, message, failedAt = models.StatusTimeLimitExceeded,
				fmt.Sprintf("exceeded %s on test %d", executionTimeout, rowCount), rowCount
		case res.ExitCode != 0:
			status, message, failedAt = models.StatusRuntimeError, truncate(res.Stderr, 400), rowCount
		case strings.TrimSpace(res.Stdout) != strings.TrimSpace(expected):
			status, message, failedAt = models.StatusWrongAnswer,
				fmt.Sprintf("expected %q, got %q",
					truncate(strings.TrimSpace(expected), 120),
					truncate(strings.TrimSpace(res.Stdout), 120)), rowCount
		default:
			continue
		}
		break
	}
	if err := rows.Err(); err != nil {
		log.Printf("RunCode iter: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to read test cases")
		return
	}
	if rowCount == 0 {
		status, message = models.StatusNoTestCases, "this problem has no test cases yet"
	}

	subID, err := recordSubmission(req.ID, c.UserID, language, req.Code, status,
		totalRuntime, failedAt, message)
	if err != nil {
		log.Printf("RunCode record submission: %v", err)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"submission_id": subID,
		"status":        status,
		"success":       status == models.StatusAccepted,
		"totalRuntime":  totalRuntime,
		"failedAt":      failedAt,
		"message":       message,
		"language":      language,
	})
}

// CustomRunCode runs the code once with the user-provided stdin and returns
// the captured output (Playground endpoint).
func CustomRunCode(w http.ResponseWriter, r *http.Request) {
	var req models.CustomCodeData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	language, _ := Language(req.Language)

	ctx, cancel := context.WithTimeout(r.Context(), executionTimeout)
	defer cancel()

	res, compileLog, compileFail, err := runUserCode(ctx, language, req.Code, req.Input)
	if err != nil {
		log.Printf("CustomRunCode exec: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to execute code")
		return
	}

	var output string
	switch {
	case compileFail:
		output = "[compilation error]\n" + compileLog
	case res.TimedOut:
		output = strings.TrimRight(res.Stdout, "\n") +
			"\n[time limit exceeded after " +
			fmt.Sprintf("%.1fs", executionTimeout.Seconds()) + "]"
	case res.Stderr != "":
		output = strings.TrimRight(res.Stdout, "\n") + "\n" + strings.TrimSpace(res.Stderr)
	default:
		output = strings.TrimRight(res.Stdout, "\n")
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"output":   strings.TrimSpace(output),
		"runtime":  res.Runtime,
		"language": language,
	})
}

// RunSampleCode runs the code against only the sample test cases and returns
// per-case results. Does not record a submission. Auth optional.
func RunSampleCode(w http.ResponseWriter, r *http.Request) {
	var req models.CodeData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	language, _ := Language(req.Language)

	rows, err := DB().Query(
		`SELECT input, output FROM testcases WHERE id=$1 AND sample=true`, req.ID)
	if err != nil {
		log.Printf("RunSampleCode query: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch sample test cases")
		return
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	allPassed := true
	rowCount := 0
	var compileFailOnce bool
	var compileLogOnce string

	for rows.Next() {
		rowCount++
		var input, expected string
		if err := rows.Scan(&input, &expected); err != nil {
			log.Printf("RunSampleCode scan: %v", err)
			httpError(w, http.StatusInternalServerError, "failed to scan test case")
			return
		}

		// If we already saw a compile error on the first test, surface it
		// without re-running. (Compile failures are file-level, not per-input.)
		if compileFailOnce {
			results = append(results, map[string]interface{}{
				"input":    input,
				"expected": expected,
				"output":   "compilation error: " + compileLogOnce,
				"result":   false,
				"runtime":  "-",
			})
			allPassed = false
			continue
		}

		ctx, cancel := context.WithTimeout(r.Context(), executionTimeout)
		res, compileLog, compileFail, runErr := runUserCode(ctx, language, req.Code, input)
		cancel()

		actual := strings.TrimRight(res.Stdout, "\n")
		passed := false
		switch {
		case compileFail:
			compileFailOnce = true
			compileLogOnce = truncate(compileLog, 300)
			actual = "compilation error: " + compileLogOnce
		case runErr != nil:
			actual = "execution error: " + runErr.Error()
		case res.TimedOut:
			actual = "time limit exceeded"
		case res.ExitCode != 0:
			actual = strings.TrimSpace(res.Stderr)
		default:
			passed = strings.TrimSpace(actual) == strings.TrimSpace(expected)
		}
		if !passed {
			allPassed = false
		}

		results = append(results, map[string]interface{}{
			"input":    input,
			"expected": expected,
			"output":   actual,
			"result":   passed,
			"runtime":  fmt.Sprintf("%.3fs", res.Runtime),
		})
	}
	if err := rows.Err(); err != nil {
		log.Printf("RunSampleCode iter: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to read test cases")
		return
	}
	if rowCount == 0 {
		allPassed = false
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":  allPassed,
		"results":  results,
		"language": language,
	})
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
