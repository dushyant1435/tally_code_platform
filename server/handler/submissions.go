package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"y/models"

	"github.com/gorilla/mux"
)

// recordSubmission inserts a row into the `submissions` table. Used by
// RunCode (full judging) — see runCode.go. Returns the new submission id.
func recordSubmission(
	problemID, userID int,
	language string,
	code string,
	status models.SubmissionStatus,
	runtime float64,
	failedTest int,
	message string,
) (int, error) {
	var (
		id          int
		runtimePtr  interface{} = runtime
		failedPtr   interface{}
		messagePtr  interface{}
	)
	if failedTest > 0 {
		failedPtr = failedTest
	}
	if message != "" {
		messagePtr = message
	}

	err := DB().QueryRow(
		`INSERT INTO submissions
		 (problem_id, user_id, language, code, status, runtime_seconds, failed_test, message)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		problemID, userID, language, code, string(status),
		runtimePtr, failedPtr, messagePtr,
	).Scan(&id)
	return id, err
}

// ListMySubmissions returns the current user's submissions, optionally filtered by problem.
// GET /api/v1/submissions?problem_id=<id>
func ListMySubmissions(w http.ResponseWriter, r *http.Request) {
	c := ClaimsFromContext(r.Context())
	if c == nil {
		httpError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	args := []interface{}{c.UserID}
	q := `SELECT s.id, s.problem_id, p.name, s.user_id, u.username,
	             s.language, s.status, s.runtime_seconds, s.failed_test, s.message, s.created_at
	      FROM submissions s
	      JOIN problems p ON p.id = s.problem_id
	      JOIN users    u ON u.id = s.user_id
	      WHERE s.user_id = $1`

	if pid := r.URL.Query().Get("problem_id"); pid != "" {
		if _, err := strconv.Atoi(pid); err == nil {
			q += ` AND s.problem_id = $2`
			args = append(args, pid)
		}
	}
	q += ` ORDER BY s.created_at DESC LIMIT 100`

	rows, err := DB().Query(q, args...)
	if err != nil {
		log.Printf("ListMySubmissions: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch submissions")
		return
	}
	defer rows.Close()

	writeSubmissions(w, rows)
}

// GetSubmission returns one submission (including its code) — only the owner
// or an admin may view a submission's code.
func GetSubmission(w http.ResponseWriter, r *http.Request) {
	c := ClaimsFromContext(r.Context())
	if c == nil {
		httpError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		httpError(w, http.StatusBadRequest, "invalid submission id")
		return
	}

	var s models.Submission
	err = DB().QueryRow(
		`SELECT s.id, s.problem_id, p.name, s.user_id, u.username, s.language, s.code,
		        s.status, s.runtime_seconds, s.failed_test, s.message, s.created_at
		 FROM submissions s
		 JOIN problems p ON p.id = s.problem_id
		 JOIN users    u ON u.id = s.user_id
		 WHERE s.id = $1`, id,
	).Scan(
		&s.ID, &s.ProblemID, &s.ProblemName, &s.UserID, &s.Username,
		&s.Language, &s.Code, &s.Status,
		&s.RuntimeSeconds, &s.FailedTest, &s.Message, &s.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		httpError(w, http.StatusNotFound, "submission not found")
		return
	}
	if err != nil {
		log.Printf("GetSubmission: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch submission")
		return
	}
	if s.UserID != c.UserID && c.Role != string(models.RoleAdmin) {
		httpError(w, http.StatusForbidden, "not allowed to view this submission")
		return
	}
	writeJSON(w, http.StatusOK, s)
}

// ListProblemSubmissions returns the current user's submissions for one problem.
// GET /api/v1/problem/{id}/submissions
func ListProblemSubmissions(w http.ResponseWriter, r *http.Request) {
	c := ClaimsFromContext(r.Context())
	if c == nil {
		httpError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	pid, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		httpError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	rows, err := DB().Query(
		`SELECT s.id, s.problem_id, p.name, s.user_id, u.username,
		        s.language, s.status, s.runtime_seconds, s.failed_test, s.message, s.created_at
		 FROM submissions s
		 JOIN problems p ON p.id = s.problem_id
		 JOIN users    u ON u.id = s.user_id
		 WHERE s.user_id = $1 AND s.problem_id = $2
		 ORDER BY s.created_at DESC LIMIT 50`,
		c.UserID, pid,
	)
	if err != nil {
		log.Printf("ListProblemSubmissions: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch submissions")
		return
	}
	defer rows.Close()

	writeSubmissions(w, rows)
}

// AdminListAllSubmissions lists every submission across every user — admin only.
func AdminListAllSubmissions(w http.ResponseWriter, r *http.Request) {
	rows, err := DB().Query(
		`SELECT s.id, s.problem_id, p.name, s.user_id, u.username,
		        s.language, s.status, s.runtime_seconds, s.failed_test, s.message, s.created_at
		 FROM submissions s
		 JOIN problems p ON p.id = s.problem_id
		 JOIN users    u ON u.id = s.user_id
		 ORDER BY s.created_at DESC LIMIT 200`,
	)
	if err != nil {
		log.Printf("AdminListAllSubmissions: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch submissions")
		return
	}
	defer rows.Close()

	writeSubmissions(w, rows)
}

func writeSubmissions(w http.ResponseWriter, rows *sql.Rows) {
	out := make([]models.Submission, 0)
	for rows.Next() {
		var s models.Submission
		if err := rows.Scan(
			&s.ID, &s.ProblemID, &s.ProblemName, &s.UserID, &s.Username,
			&s.Language, &s.Status,
			&s.RuntimeSeconds, &s.FailedTest, &s.Message, &s.CreatedAt,
		); err != nil {
			log.Printf("writeSubmissions scan: %v", err)
			httpError(w, http.StatusInternalServerError, "failed to read submissions")
			return
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		log.Printf("writeSubmissions iter: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to read submissions")
		return
	}
	writeJSON(w, http.StatusOK, out)
}
