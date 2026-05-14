package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"y/models"
)

const (
	executionTimeout = 5 * time.Second
	pythonBinary     = "python3"
)

type execResult struct {
	Stdout   string
	Stderr   string
	Runtime  float64
	TimedOut bool
	ExitCode int
}

// runPython writes code to a fresh temp file and executes it once with the
// supplied stdin, returning captured stdout/stderr and elapsed time.
func runPython(ctx context.Context, code, stdin string) (execResult, error) {
	f, err := os.CreateTemp("", "tally-*.py")
	if err != nil {
		return execResult{}, fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(code); err != nil {
		f.Close()
		return execResult{}, fmt.Errorf("write temp file: %w", err)
	}
	if err := f.Close(); err != nil {
		return execResult{}, fmt.Errorf("close temp file: %w", err)
	}

	cmd := exec.CommandContext(ctx, pythonBinary, f.Name())
	cmd.Stdin = strings.NewReader(stdin)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Run()
	elapsed := time.Since(start).Seconds()

	res := execResult{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Runtime: elapsed,
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		res.TimedOut = true
		return res, nil
	}
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			res.ExitCode = exitErr.ExitCode()
			return res, nil
		}
		return res, err
	}
	return res, nil
}

// RunCode executes the submitted code against every test case for a problem.
// Records a submissions row with the resulting verdict and returns it.
// Requires authentication.
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
	if req.Language == "" {
		req.Language = "python"
	}

	db := DB()
	rows, err := db.Query(`SELECT input, output FROM testcases WHERE id=$1`, req.ID)
	if err != nil {
		log.Printf("RunCode query: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch test cases")
		return
	}
	defer rows.Close()

	var (
		status         = models.StatusAccepted
		message        string
		failedAt       int
		totalRuntime   float64
		rowCount       int
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
		res, runErr := runPython(ctx, req.Code, input)
		cancel()

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
			snippet := strings.TrimSpace(res.Stderr)
			if len(snippet) > 200 {
				snippet = snippet[:200] + "..."
			}
			status, message, failedAt = models.StatusRuntimeError, snippet, rowCount
		case strings.TrimSpace(res.Stdout) != strings.TrimSpace(expected):
			status, message, failedAt = models.StatusWrongAnswer,
				fmt.Sprintf("expected %q, got %q",
					strings.TrimSpace(expected),
					strings.TrimSpace(res.Stdout)), rowCount
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

	subID, err := recordSubmission(req.ID, c.UserID, req.Language, req.Code, status,
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
	})
}

// CustomRunCode runs the code once with the user-provided stdin and returns
// the captured output (Playground endpoint). Does not require auth.
func CustomRunCode(w http.ResponseWriter, r *http.Request) {
	var req models.CustomCodeData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), executionTimeout)
	defer cancel()

	res, err := runPython(ctx, req.Code, req.Input)
	if err != nil {
		log.Printf("CustomRunCode exec: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to execute code")
		return
	}

	output := strings.TrimRight(res.Stdout, "\n")
	switch {
	case res.TimedOut:
		output = strings.TrimSpace(output + "\n[time limit exceeded after " +
			fmt.Sprintf("%.1fs", executionTimeout.Seconds()) + "]")
	case res.Stderr != "":
		output = strings.TrimSpace(output + "\n" + res.Stderr)
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"output":  output,
		"runtime": res.Runtime,
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

	for rows.Next() {
		rowCount++
		var input, expected string
		if err := rows.Scan(&input, &expected); err != nil {
			log.Printf("RunSampleCode scan: %v", err)
			httpError(w, http.StatusInternalServerError, "failed to scan test case")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), executionTimeout)
		res, runErr := runPython(ctx, req.Code, input)
		cancel()

		actual := strings.TrimRight(res.Stdout, "\n")
		passed := false
		switch {
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
		"success": allPassed,
		"results": results,
	})
}
