package models

import "time"

// ---------- users ---------------------------------------------------------

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	// Allow either email or username in the same field.
	UsernameOrEmail string `json:"username"`
	Password        string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ---------- problems ------------------------------------------------------

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type Problem struct {
	ID           int        `json:"id"`
	UserId       int        `json:"user_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Constraints  *string    `json:"constraints,omitempty"`
	InputFormat  *string    `json:"input_format,omitempty"`
	OutputFormat *string    `json:"output_format,omitempty"`
	Difficulty   Difficulty `json:"difficulty"`
	Tags         []string   `json:"tags"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ---------- test cases ----------------------------------------------------

type TestCase struct {
	ID     int    `json:"id"`
	Input  string `json:"input"`
	Output string `json:"output"`
	Sample bool   `json:"sample"`
}

// ---------- code-run payloads ---------------------------------------------

type CodeData struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Language string `json:"language"`
}

type CustomCodeData struct {
	Code     string `json:"code"`
	Input    string `json:"input"`
	Language string `json:"language"`
}

// ---------- submissions ---------------------------------------------------

type SubmissionStatus string

const (
	StatusAccepted          SubmissionStatus = "accepted"
	StatusWrongAnswer       SubmissionStatus = "wrong_answer"
	StatusTimeLimitExceeded SubmissionStatus = "time_limit_exceeded"
	StatusRuntimeError      SubmissionStatus = "runtime_error"
	StatusCompileError      SubmissionStatus = "compilation_error"
	StatusNoTestCases       SubmissionStatus = "no_test_cases"
	StatusServerError       SubmissionStatus = "server_error"
)

type Submission struct {
	ID             int              `json:"id"`
	ProblemID      int              `json:"problem_id"`
	ProblemName    string           `json:"problem_name,omitempty"`
	UserID         int              `json:"user_id"`
	Username       string           `json:"username,omitempty"`
	Language       string           `json:"language"`
	Code           string           `json:"code,omitempty"`
	Status         SubmissionStatus `json:"status"`
	RuntimeSeconds *float64         `json:"runtime_seconds,omitempty"`
	FailedTest     *int             `json:"failed_test,omitempty"`
	Message        *string          `json:"message,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
}
