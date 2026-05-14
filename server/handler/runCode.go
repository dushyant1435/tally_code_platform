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
	Stdout  string
	Stderr  string
	Runtime float64
	TimedOut bool
}

// runPython writes code to a fresh temp file and executes it once with the
// supplied stdin, returning the captured stdout/stderr and elapsed time.
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
		// Non-zero exit is not fatal for us; we still want stdout/stderr.
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return res, nil
		}
		return res, err
	}
	return res, nil
}

// RunCode executes the submitted code against every test case for a problem.
// On full pass it records a row in `submission` (idempotent).
func RunCode(w http.ResponseWriter, r *http.Request) {
	var req models.CodeData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	db := DB()
	rows, err := db.Query(`SELECT input, output FROM testcases WHERE id=$1`, req.ID)
	if err != nil {
		log.Printf("RunCode query: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch test cases")
		return
	}
	defer rows.Close()

	allPassed := true
	totalRuntime := 0.0
	rowCount := 0
	var failedAt int
	var failMessage string

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

		if runErr != nil {
			log.Printf("RunCode exec: %v", runErr)
			allPassed = false
			failedAt = rowCount
			failMessage = "execution error"
			break
		}
		if res.TimedOut {
			allPassed = false
			failedAt = rowCount
			failMessage = "time limit exceeded"
			break
		}

		if res.Runtime > totalRuntime {
			totalRuntime = res.Runtime
		}

		if strings.TrimSpace(res.Stdout) != strings.TrimSpace(expected) {
			allPassed = false
			failedAt = rowCount
			failMessage = "wrong answer"
			break
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("RunCode iter: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to read test cases")
		return
	}

	if rowCount == 0 {
		allPassed = false
		failMessage = "no test cases configured for this problem"
	}

	if allPassed {
		_, err := db.Exec(
			`INSERT INTO submission (id, user_id) VALUES ($1, $2)
			 ON CONFLICT (id, user_id) DO NOTHING`,
			req.ID, req.UserID,
		)
		if err != nil {
			log.Printf("RunCode submission insert: %v", err)
		}
	}

	resp := map[string]interface{}{
		"success":      allPassed,
		"totalRuntime": totalRuntime,
		"memoryUsed":   0,
		"failedAt":     failedAt,
		"message":      failMessage,
	}
	writeJSON(w, http.StatusOK, resp)
}

// CustomRunCode runs the code once with the user-provided stdin and returns
// the captured output (Playground endpoint).
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
	if res.TimedOut {
		output = strings.TrimSpace(output + "\n[time limit exceeded after " + fmt.Sprintf("%.1fs", executionTimeout.Seconds()) + "]")
	} else if res.Stderr != "" {
		output = strings.TrimSpace(output + "\n" + res.Stderr)
	}

	writeJSON(w, http.StatusOK, map[string]string{"output": output})
}

// RunSampleCode runs the code against only the sample test cases and returns
// per-case results so the UI can show actual-vs-expected.
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
		case res.Stderr != "":
			actual = strings.TrimSpace(res.Stderr)
		default:
			passed = strings.TrimSpace(actual) == strings.TrimSpace(expected)
		}
		if !passed {
			allPassed = false
		}

		results = append(results, map[string]interface{}{
			"input":       input,
			"expected":    expected,
			"output":      actual,
			"result":      passed,
			"runtime":     fmt.Sprintf("%.3fs", res.Runtime),
			"memory_used": "-",
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
