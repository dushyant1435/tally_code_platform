package router

import (
	"net/http"
	"os"
	"strings"
	handler "y/handler"

	"github.com/gorilla/mux"
)

// corsMiddleware adds permissive CORS headers and short-circuits OPTIONS
// preflight requests so the React dev server can talk to us.
func corsMiddleware(next http.Handler) http.Handler {
	allowed := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowed == "" {
		allowed = "*"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowed == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" && contains(strings.Split(allowed, ","), origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if strings.TrimSpace(v) == s {
			return true
		}
	}
	return false
}

// Router builds the application router used by main.go.
func Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	r.HandleFunc("/api/v1/problem/{id}", handler.GetProblem).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/problems", handler.GetAllProblems).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/newproblem", handler.CreateProblem).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/testcases/{id}", handler.GetTestCasesByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/problem/{id}/sampleTestCases", handler.GetSampleTestCasesByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/createTestCase", handler.CreateTestCase).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/runCode", handler.RunCode).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/runSampleCode", handler.RunSampleCode).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/runCustomCode", handler.CustomRunCode).Methods("POST", "OPTIONS")

	return corsMiddleware(r)
}
