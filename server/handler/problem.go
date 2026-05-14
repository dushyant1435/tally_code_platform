package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"y/models"

	"github.com/gorilla/mux"
)

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write json: %v", err)
	}
}

func httpError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func GetProblemByID(id int64) (models.Problem, error) {
	var problem models.Problem
	row := DB().QueryRow(
		`SELECT id, user_id, name, description, constraints, input_format, output_format
		 FROM problems WHERE id=$1`, id)

	err := row.Scan(
		&problem.ID,
		&problem.UserId,
		&problem.Name,
		&problem.Description,
		&problem.Constraints,
		&problem.InputFormat,
		&problem.OutputFormat,
	)
	return problem, err
}

// GetProblem returns a single problem by ID.
func GetProblem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		httpError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	problem, err := GetProblemByID(int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		httpError(w, http.StatusNotFound, "problem not found")
		return
	}
	if err != nil {
		log.Printf("GetProblem scan: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch problem")
		return
	}
	writeJSON(w, http.StatusOK, problem)
}

// GetAllProblems returns every problem, plus a "status" flag indicating whether
// the given user has already solved it. user_id is taken from the query string,
// e.g. /api/v1/problems?user_id=123.
func GetAllProblems(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		httpError(w, http.StatusBadRequest, "missing user_id query parameter")
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httpError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	db := DB()
	rows, err := db.Query(
		`SELECT p.id, p.user_id, p.name, p.description, p.constraints,
		        p.input_format, p.output_format,
		        EXISTS (SELECT 1 FROM submission s WHERE s.id = p.id AND s.user_id = $1) AS solved
		 FROM problems p ORDER BY p.id`, userID)
	if err != nil {
		log.Printf("GetAllProblems query: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch problems")
		return
	}
	defer rows.Close()

	problems := make([]map[string]interface{}, 0)
	for rows.Next() {
		var p models.Problem
		var solved bool
		if err := rows.Scan(
			&p.ID, &p.UserId, &p.Name, &p.Description,
			&p.Constraints, &p.InputFormat, &p.OutputFormat, &solved,
		); err != nil {
			log.Printf("GetAllProblems scan: %v", err)
			httpError(w, http.StatusInternalServerError, "failed to scan problem")
			return
		}
		problems = append(problems, map[string]interface{}{
			"id":            p.ID,
			"user_id":       p.UserId,
			"name":          p.Name,
			"description":   p.Description,
			"constraints":   p.Constraints,
			"input_format":  p.InputFormat,
			"output_format": p.OutputFormat,
			"status":        solved,
		})
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetAllProblems iter: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to read problems")
		return
	}
	writeJSON(w, http.StatusOK, problems)
}

// CreateProblem inserts a new problem and returns the row with the generated ID.
func CreateProblem(w http.ResponseWriter, r *http.Request) {
	var problem models.Problem
	if err := json.NewDecoder(r.Body).Decode(&problem); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if problem.Name == "" || problem.Description == "" {
		httpError(w, http.StatusBadRequest, "name and description are required")
		return
	}

	var id int
	err := DB().QueryRow(
		`INSERT INTO problems (user_id, name, description, constraints, input_format, output_format)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		problem.UserId, problem.Name, problem.Description,
		problem.Constraints, problem.InputFormat, problem.OutputFormat,
	).Scan(&id)
	if err != nil {
		log.Printf("CreateProblem insert: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to create problem")
		return
	}

	problem.ID = id
	writeJSON(w, http.StatusCreated, problem)
}
