package router

import (
	handler "y/handler"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/problem/{id}", handler.GetProblem).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/problems", handler.GetAllProblems).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/newproblem", handler.CreateProblem).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/testcases/{id}", handler.GetTestCasesByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/problem/{id}/sampleTestCases", handler.GetSampleTestCasesByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/createTestCase", handler.CreateTestCase).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/runCode", handler.RunCode).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/runSampleCode", handler.RunSampleCode).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/runCustomCode", handler.CustomRunCode).Methods("POST", "OPTIONS")
	return router
}
