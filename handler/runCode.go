package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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
	var totalRuntime float64
	var maxMemoryUsed int64

	for rows.Next() {
		var input, expectedOutput string
		if err := rows.Scan(&input, &expectedOutput); err != nil {
			http.Error(w, fmt.Sprintf("Unable to scan row: %v", err), http.StatusInternalServerError)
			return
		}

		// Use 'time' command to measure execution time and memory usage
		cmd := exec.Command("bash", "-c", fmt.Sprintf("/usr/bin/time -f '%%e %%M' python3 %s", filename))
		cmd.Stdin = strings.NewReader(input)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing code: %v", err)
			allTestsPassed = false
			break
		}

		// The output from the 'time' command contains time and memory usage
		outputParts := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(outputParts) < 2 {
			http.Error(w, "Failed to get execution metrics", http.StatusInternalServerError)
			return
		}
		metrics := strings.Fields(outputParts[len(outputParts)-1])
		if len(metrics) < 2 {
			http.Error(w, "Failed to parse execution metrics", http.StatusInternalServerError)
			return
		}

		// Convert the metrics to appropriate types
		runtime, err := strconv.ParseFloat(metrics[0], 64)
		if err != nil {
			http.Error(w, "Failed to parse runtime", http.StatusInternalServerError)
			return
		}
		memoryUsed, err := strconv.ParseInt(metrics[1], 10, 64)
		if err != nil {
			http.Error(w, "Failed to parse memory usage", http.StatusInternalServerError)
			return
		}

		// Accumulate total runtime and track the maximum memory used
		if totalRuntime <= runtime {
			totalRuntime = runtime
		}
		if memoryUsed > maxMemoryUsed {
			maxMemoryUsed = memoryUsed
		}

		// Compare the output with the expected output
		if strings.TrimSpace(outputParts[0]) != strings.TrimSpace(expectedOutput) {
			allTestsPassed = false
			break
		}
	}
	if allTestsPassed {
		insertStatement := `INSERT INTO submission (id, user_id) VALUES ($1, $2)`
		_, err = db.Exec(insertStatement, requestData.ID, requestData.UserID)
		if err != nil {
			log.Printf("Error inserting into submission table: %v", err)
		}
	}

	// Return the overall success status, total runtime, and max memory usage
	response := map[string]interface{}{
		"success":      allTestsPassed,
		"totalRuntime": totalRuntime,
		"memoryUsed":   maxMemoryUsed, // in kilobytes
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

	// Generate a unique filename based on problem ID
	filename := fmt.Sprintf("temp_%d.py", requestData.ID)

	// Prepare the Python code with metrics collection
	codeWithMetrics := `
import sys
import time
import tracemalloc

tracemalloc.start()
start_time = time.time()

` + requestData.Code + `

end_time = time.time()
current, peak = tracemalloc.get_traced_memory()
execution_time = end_time - start_time

sys.stderr.write(f"Execution Time: {execution_time} seconds\n")
sys.stderr.write(f"Peak Memory Usage: {peak / 1024} KB\n")
tracemalloc.stop()
`

	// Write the modified code to the temporary file
	err = os.WriteFile(filename, []byte(codeWithMetrics), 0644)
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
		outputStr := string(output)

		// Parse the execution time and memory usage from stderr
		lines := strings.Split(outputStr, "\n")
		runtime := ""
		memoryUsed := ""
		if len(lines) > 2 {
			runtime = strings.TrimPrefix(lines[0], "Execution Time: ")
			memoryUsed = strings.TrimPrefix(lines[1], "Peak Memory Usage: ")
		}

		if err != nil {
			log.Printf("Error executing code: %v", err)
			allTestsPassed = false
			results = append(results, map[string]interface{}{
				"input":       input,
				"expected":    expectedOutput,
				"output":      "error",
				"result":      false,
				"runtime":     runtime,
				"memory_used": memoryUsed,
			})
			continue
		}

		// Compare the output with the expected output
		testPassed := lines[2] == strings.TrimSpace(expectedOutput)
		if !testPassed {
			allTestsPassed = false
		}

		results = append(results, map[string]interface{}{
			"input":       input,
			"expected":    expectedOutput,
			"output":      lines[2],
			"result":      testPassed,
			"runtime":     runtime,
			"memory_used": memoryUsed,
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
