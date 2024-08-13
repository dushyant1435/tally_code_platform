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
	// Parse the JSON request body
	var testcase models.TestCase
	err := json.NewDecoder(r.Body).Decode(&testcase)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create a connection to the database
	db := createConnection()

	// Close the db connection
	defer db.Close()

	// Insert query
	sqlStatement := `
    INSERT INTO testcases (id, input, output, sample)
    VALUES ($1, $2, $3, $4)`

	// Execute the SQL statement
	_, err = db.Exec(sqlStatement, testcase.ID, testcase.Input, testcase.Output, testcase.Sample)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with status created
	w.WriteHeader(http.StatusCreated)
}

func GetTestCasesByID(w http.ResponseWriter, r *http.Request) {
	// Get the id from the request parameters
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the database connection
	db := createConnection()

	// Close the database connection
	defer db.Close()

	// Define the query to get test cases by id
	sqlStatement := `SELECT id, input, output, sample FROM testcases WHERE id=$1`

	// Execute the query
	rows, err := db.Query(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	// Create a slice to store the test cases
	var testCases []models.TestCase

	// Iterate over the rows and add to the slice
	for rows.Next() {
		var testCase models.TestCase

		err = rows.Scan(&testCase.ID, &testCase.Input, &testCase.Output, &testCase.Sample)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		testCases = append(testCases, testCase)
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		log.Fatalf("Error during iteration over rows. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testCases)
}

func GetSampleTestCasesByID(w http.ResponseWriter, r *http.Request) {
	// Get the id from the request parameters
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the database connection
	db := createConnection()

	// Close the database connection
	defer db.Close()

	// Define the query to get test cases by id
	sqlStatement := `SELECT id, input, output, sample FROM testcases WHERE id=$1 AND sample=true`

	// Execute the query
	rows, err := db.Query(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	// Create a slice to store the test cases
	var testCases []models.TestCase

	// Iterate over the rows and add to the slice
	for rows.Next() {
		var testCase models.TestCase

		err = rows.Scan(&testCase.ID, &testCase.Input, &testCase.Output, &testCase.Sample)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		testCases = append(testCases, testCase)
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		log.Fatalf("Error during iteration over rows. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testCases)
}
