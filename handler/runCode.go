package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"y/models"
)

func RunCode(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var requestData models.CodeData

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to decode the request body: %v", err), http.StatusBadRequest)
		return
	}

	db := createConnection()
	defer db.Close()

	// Fetch test cases using the problem ID
	sqlStatement := `SELECT input, output FROM testcases WHERE id=$1 AND sample=true`
	rows, err := db.Query(sqlStatement, requestData.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to fetch test cases: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Generate a unique filename based on user ID and problem ID
	filename := fmt.Sprintf("temp_%d_%d.py", requestData.UserID, requestData.ID)

	// Write the code to the temporary file
	err = os.WriteFile(filename, []byte(requestData.Code), 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to write temporary Python file: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(filename) // Ensure the file is removed after execution

	// Execute the code for each test case and compare results
	allTestsPassed := true

	for rows.Next() {
		var input, expectedOutput string
		if err := rows.Scan(&input, &expectedOutput); err != nil {
			http.Error(w, fmt.Sprintf("Unable to scan row: %v", err), http.StatusInternalServerError)
			return
		}

		// Execute the Python file with the input as argument
		cmd := exec.Command("python3", filename)
		cmd.Stdin = strings.NewReader(input)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing code: %v", err)
			allTestsPassed = false
			break
		}

		// Compare the output with the expected output
		if strings.TrimSpace(string(output)) != strings.TrimSpace(expectedOutput) {
			allTestsPassed = false
			break
		}
	}

	// Return true if all test cases passed, false otherwise
	response := map[string]bool{"success": allTestsPassed}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func RunSampleCode(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var requestData models.CodeData

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to decode the request body: %v", err), http.StatusBadRequest)
		return
	}

	db := createConnection()
	defer db.Close()

	// Fetch sample test cases using the problem ID
	sqlStatement := `SELECT input, output FROM testcases WHERE id=$1 AND sample=true`
	rows, err := db.Query(sqlStatement, requestData.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to fetch test cases: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Generate a unique filename based on user ID and problem ID
	filename := fmt.Sprintf("temp_%d_%d.py", requestData.UserID, requestData.ID)

	// Write the code to the temporary file
	err = os.WriteFile(filename, []byte(requestData.Code), 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to write temporary Python file: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(filename) // Ensure the file is removed after execution

	// Execute the code for each test case and compare results
	allTestsPassed := true
	var results []map[string]interface{}

	for rows.Next() {
		var input, expectedOutput string
		if err := rows.Scan(&input, &expectedOutput); err != nil {
			http.Error(w, fmt.Sprintf("Unable to scan row: %v", err), http.StatusInternalServerError)
			return
		}

		// Execute the Python file with the input as argument
		cmd := exec.Command("python3", filename)
		cmd.Stdin = strings.NewReader(input)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing code: %v", err)
			allTestsPassed = false
			results = append(results, map[string]interface{}{
				"input":    input,
				"expected": expectedOutput,
				"output":   "error",
				"result":   false,
			})
			continue
		}

		// Compare the output with the expected output
		testPassed := strings.TrimSpace(string(output)) == strings.TrimSpace(expectedOutput)
		if !testPassed {
			allTestsPassed = false
		}

		results = append(results, map[string]interface{}{
			"input":    input,
			"expected": expectedOutput,
			"output":   strings.TrimSpace(string(output)),
			"result":   testPassed,
		})
	}

	// Return results with the overall success status
	response := map[string]interface{}{
		"success": allTestsPassed,
		"results": results,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func CustomRunCode(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into CustomCodeData struct
	var requestData models.CustomCodeData

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to decode the request body: %v", err), http.StatusBadRequest)
		return
	}

	// Generate a unique filename
	filename := "temp_code.py"

	// Write the code to the temporary file
	err = os.WriteFile(filename, []byte(requestData.Code), 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to write temporary Python file: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(filename) // Ensure the file is removed after execution

	// Execute the code with the provided input
	cmd := exec.Command("python3", filename)
	cmd.Stdin = strings.NewReader(requestData.Input)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing code: %v", err)
		http.Error(w, fmt.Sprintf("Error executing code: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the output
	response := map[string]string{"output": strings.TrimSpace(string(output))}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
