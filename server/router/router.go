package router

import (
	"net/http"
	"os"
	"strings"
	handler "y/handler"

	"github.com/gorilla/mux"
)

// corsMiddleware adds CORS headers + handles preflight.
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

func Router() http.Handler {
	r := mux.NewRouter()

	// Liveness probe.
	r.HandleFunc("/api/v1/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// ----- Auth ----------------------------------------------------------
	r.HandleFunc("/api/v1/auth/signup", handler.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/auth/login", handler.Login).Methods("POST", "OPTIONS")
	r.Handle("/api/v1/auth/me", handler.RequireAuth(http.HandlerFunc(handler.Me))).
		Methods("GET", "OPTIONS")

	// ----- Problems (public reads with optional auth for "solved" flag) ---
	r.Handle("/api/v1/problems",
		handler.OptionalAuth(http.HandlerFunc(handler.GetAllProblems))).
		Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/problem/{id:[0-9]+}", handler.GetProblem).
		Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/testcases/{id:[0-9]+}", handler.GetTestCasesByID).
		Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/problem/{id:[0-9]+}/sampleTestCases",
		handler.GetSampleTestCasesByID).Methods("GET", "OPTIONS")

	// ----- Authoring (admin only) ---------------------------------------
	r.Handle("/api/v1/newproblem",
		handler.RequireAdmin(http.HandlerFunc(handler.CreateProblem))).
		Methods("POST", "OPTIONS")
	r.Handle("/api/v1/createTestCase",
		handler.RequireAdmin(http.HandlerFunc(handler.CreateTestCase))).
		Methods("POST", "OPTIONS")

	// ----- Code execution -----------------------------------------------
	// runCode: must be authenticated (records a submission against the user).
	r.Handle("/api/v1/runCode",
		handler.RequireAuth(http.HandlerFunc(handler.RunCode))).
		Methods("POST", "OPTIONS")
	// sample/custom: open to anyone for the Playground.
	r.HandleFunc("/api/v1/runSampleCode", handler.RunSampleCode).
		Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/runCustomCode", handler.CustomRunCode).
		Methods("POST", "OPTIONS")

	// ----- Submissions --------------------------------------------------
	r.Handle("/api/v1/submissions",
		handler.RequireAuth(http.HandlerFunc(handler.ListMySubmissions))).
		Methods("GET", "OPTIONS")
	r.Handle("/api/v1/submissions/{id:[0-9]+}",
		handler.RequireAuth(http.HandlerFunc(handler.GetSubmission))).
		Methods("GET", "OPTIONS")
	r.Handle("/api/v1/problem/{id:[0-9]+}/submissions",
		handler.RequireAuth(http.HandlerFunc(handler.ListProblemSubmissions))).
		Methods("GET", "OPTIONS")

	// ----- Admin --------------------------------------------------------
	r.Handle("/api/v1/admin/submissions",
		handler.RequireAdmin(http.HandlerFunc(handler.AdminListAllSubmissions))).
		Methods("GET", "OPTIONS")

	return corsMiddleware(r)
}
