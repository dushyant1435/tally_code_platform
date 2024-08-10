package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"y/models"

	"github.com/gorilla/mux"
)

func GetProblemByID(id int64) (models.Problem, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()
	// Create a problem of models.Problem type
	var problem models.Problem

	// Create the select SQL query
	sqlStatement := `SELECT * FROM Problems WHERE id=$1`

	// Execute the SQL statement
	row := db.QueryRow(sqlStatement, id)

	// Unmarshal the row object to problem
	err := row.Scan(
		&problem.ID,
		&problem.UserId,
		&problem.Name,
		&problem.Description,
		&problem.Constraints,
		&problem.InputFormat,
		&problem.OutputFormat,
	)

	switch err {
	case sql.ErrNoRows:
		// No rows returned; return the empty problem with no error
		return problem, nil
	case nil:
		// Successfully retrieved the problem
		return problem, nil
	default:
		// Error occurred while scanning the row
		log.Fatalf("Unable to scan the row. %v", err)
		return problem, err
	}
}

// GetProblem handles the request to get a problem by its ID
func GetProblem(w http.ResponseWriter, r *http.Request) {
	// Get the problem ID from the request parameters
	params := mux.Vars(r)

	// Convert the ID from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	Problem, err := GetProblemByID(int64(id))

	if err != nil {
		log.Fatalf("Unable to get problem. %v", err)
	}

	json.NewEncoder(w).Encode(Problem)
}

func GetAllProblems(w http.ResponseWriter, r *http.Request) {
	// Create a connection to the database
	db := createConnection()
	defer db.Close()

	// SQL query to get all problems
	sqlStatement := `SELECT id, user_id, name, description, constraints, input_format, output_format FROM problems`

	// Execute the query
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Slice to store the problems
	var problems []models.Problem

	// Iterate over the rows
	for rows.Next() {
		var problem models.Problem

		// Scan the row into the problem struct
		err = rows.Scan(&problem.ID, &problem.UserId, &problem.Name, &problem.Description, &problem.Constraints, &problem.InputFormat, &problem.OutputFormat)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Append the problem to the slice
		problems = append(problems, problem)
	}

	// Return the list of problems as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(problems)
}

func CreateProblem(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var problem models.Problem
	err := json.NewDecoder(r.Body).Decode(&problem)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := createConnection()

	// close the db connection
	defer db.Close()

	// Insert query
	sqlStatement := `
	INSERT INTO problems (user_id, name, description, constraints, input_format, output_format)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	// Execute the SQL statement
	var id int
	err = db.QueryRow(sqlStatement, problem.UserId, problem.Name, problem.Description, problem.Constraints, problem.InputFormat, problem.OutputFormat).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the created problem with the generated ID
	problem.ID = id
	fmt.Println(problem.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(problem)
}

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
