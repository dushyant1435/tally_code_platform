package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"y/models"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
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

func validDifficulty(d models.Difficulty) bool {
	switch d {
	case models.DifficultyEasy, models.DifficultyMedium, models.DifficultyHard:
		return true
	}
	return false
}

// GetProblemByID is a small helper used by runCode.go.
func GetProblemByID(id int64) (models.Problem, error) {
	var p models.Problem
	err := DB().QueryRow(
		`SELECT id, user_id, name, description, constraints, input_format, output_format,
		        difficulty, tags, created_at
		 FROM problems WHERE id=$1`, id,
	).Scan(
		&p.ID, &p.UserId, &p.Name, &p.Description, &p.Constraints,
		&p.InputFormat, &p.OutputFormat, &p.Difficulty, pq.Array(&p.Tags), &p.CreatedAt,
	)
	return p, err
}

// GetProblem returns a single problem by ID. Public route.
func GetProblem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		httpError(w, http.StatusBadRequest, "invalid problem id")
		return
	}
	p, err := GetProblemByID(int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		httpError(w, http.StatusNotFound, "problem not found")
		return
	}
	if err != nil {
		log.Printf("GetProblem: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch problem")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// GetAllProblems returns every problem. If the request is authenticated,
// each row includes a `status` boolean indicating whether the caller has
// already solved it. Anonymous callers always see status=false.
func GetAllProblems(w http.ResponseWriter, r *http.Request) {
	userID := 0
	if c := ClaimsFromContext(r.Context()); c != nil {
		userID = c.UserID
	}

	rows, err := DB().Query(
		`SELECT p.id, p.user_id, u.username, p.name, p.description,
		        p.constraints, p.input_format, p.output_format,
		        p.difficulty, p.tags, p.created_at,
		        EXISTS (
		            SELECT 1 FROM submissions s
		            WHERE s.problem_id = p.id
		              AND s.user_id    = $1
		              AND s.status     = 'accepted'
		        ) AS solved
		 FROM problems p
		 LEFT JOIN users u ON u.id = p.user_id
		 ORDER BY p.id`, userID,
	)
	if err != nil {
		log.Printf("GetAllProblems query: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch problems")
		return
	}
	defer rows.Close()

	problems := make([]map[string]interface{}, 0)
	for rows.Next() {
		var (
			p        models.Problem
			username sql.NullString
			solved   bool
		)
		if err := rows.Scan(
			&p.ID, &p.UserId, &username, &p.Name, &p.Description,
			&p.Constraints, &p.InputFormat, &p.OutputFormat,
			&p.Difficulty, pq.Array(&p.Tags), &p.CreatedAt, &solved,
		); err != nil {
			log.Printf("GetAllProblems scan: %v", err)
			httpError(w, http.StatusInternalServerError, "failed to scan problem")
			return
		}
		problems = append(problems, map[string]interface{}{
			"id":            p.ID,
			"user_id":       p.UserId,
			"author":        username.String,
			"name":          p.Name,
			"description":   p.Description,
			"constraints":   p.Constraints,
			"input_format":  p.InputFormat,
			"output_format": p.OutputFormat,
			"difficulty":    p.Difficulty,
			"tags":          p.Tags,
			"created_at":    p.CreatedAt,
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

// CreateProblem inserts a new problem authored by the authenticated user.
// Admin-only; auth middleware enforces that on the route.
func CreateProblem(w http.ResponseWriter, r *http.Request) {
	c := ClaimsFromContext(r.Context())
	if c == nil {
		httpError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var p models.Problem
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	p.Name = strings.TrimSpace(p.Name)
	p.Description = strings.TrimSpace(p.Description)
	if p.Name == "" || p.Description == "" {
		httpError(w, http.StatusBadRequest, "name and description are required")
		return
	}
	if p.Difficulty == "" {
		p.Difficulty = models.DifficultyEasy
	}
	if !validDifficulty(p.Difficulty) {
		httpError(w, http.StatusBadRequest, "difficulty must be easy, medium or hard")
		return
	}
	if p.Tags == nil {
		p.Tags = []string{}
	}

	var id int
	err := DB().QueryRow(
		`INSERT INTO problems
		  (user_id, name, description, constraints, input_format, output_format, difficulty, tags)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		c.UserID, p.Name, p.Description, p.Constraints, p.InputFormat,
		p.OutputFormat, string(p.Difficulty), pq.Array(p.Tags),
	).Scan(&id)
	if err != nil {
		log.Printf("CreateProblem insert: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to create problem")
		return
	}
	p.ID = id
	p.UserId = c.UserID
	writeJSON(w, http.StatusCreated, p)
}
