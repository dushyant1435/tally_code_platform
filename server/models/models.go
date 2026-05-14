package models

type Problem struct {
	ID           int     `json:"id"`
	UserId       int     `json:"user_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Constraints  *string `json:"constraints,omitempty"`   // use *string to allow null values
	InputFormat  *string `json:"input_format,omitempty"`  // use *string to allow null values
	OutputFormat *string `json:"output_format,omitempty"` // use *string to allow null values
	// Status       bool  `json:"status`
}

type TestCase struct {
	ID     int    `json:"id"`
	Input  string `json:"input"`
	Output string `json:"output"`
	Sample bool   `json:"sample"`
}


type CodeData struct {
	ID     int    `json:"id"` 
	Code   string `json:"code"`
	UserID int    `json:"user_id"`
}

type CustomCodeData struct {
    Code  string `json:"code"`
    Input string `json:"input"`
}

type RequestBody struct {
	UserID int64 `json:"user_id"`
}