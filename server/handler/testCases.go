package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"y/models"

	"github.com/gorilla/mux"
)

func CreateTestCase(w http.ResponseWriter, r *http.Request) {
	var tc models.TestCase
	if err := json.NewDecoder(r.Body).Decode(&tc); err != nil {
		httpError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if tc.ID == 0 || tc.Input == "" || tc.Output == "" {
		httpError(w, http.StatusBadRequest, "id, input and output are required")
		return
	}

	_, err := DB().Exec(
		`INSERT INTO testcases (id, input, output, sample) VALUES ($1, $2, $3, $4)`,
		tc.ID, tc.Input, tc.Output, tc.Sample,
	)
	if err != nil {
		log.Printf("CreateTestCase insert: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to create test case")
		return
	}
	writeJSON(w, http.StatusCreated, tc)
}

func fetchTestCases(w http.ResponseWriter, r *http.Request, sampleOnly bool) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		httpError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	query := `SELECT id, input, output, sample FROM testcases WHERE id=$1`
	if sampleOnly {
		query += ` AND sample=true`
	}

	rows, err := DB().Query(query, id)
	if err != nil {
		log.Printf("fetchTestCases query: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to fetch test cases")
		return
	}
	defer rows.Close()

	cases := make([]models.TestCase, 0)
	for rows.Next() {
		var c models.TestCase
		if err := rows.Scan(&c.ID, &c.Input, &c.Output, &c.Sample); err != nil {
			log.Printf("fetchTestCases scan: %v", err)
			httpError(w, http.StatusInternalServerError, "failed to scan test case")
			return
		}
		cases = append(cases, c)
	}
	if err := rows.Err(); err != nil {
		log.Printf("fetchTestCases iter: %v", err)
		httpError(w, http.StatusInternalServerError, "failed to read test cases")
		return
	}
	writeJSON(w, http.StatusOK, cases)
}

func GetTestCasesByID(w http.ResponseWriter, r *http.Request) {
	fetchTestCases(w, r, false)
}

func GetSampleTestCasesByID(w http.ResponseWriter, r *http.Request) {
	fetchTestCases(w, r, true)
}
