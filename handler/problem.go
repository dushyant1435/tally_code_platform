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


