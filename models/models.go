package models

// User schema of the user table
type Stock struct {
	StockID int64  `json:"stockid"`
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Company string `json:"company"`
}

type Problem struct {
	ID           int     `json:"id"`
	UserId       int     `json:"user_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Constraints  *string `json:"constraints,omitempty"`   // use *string to allow null values
	InputFormat  *string `json:"input_format,omitempty"`  // use *string to allow null values
	OutputFormat *string `json:"output_format,omitempty"` // use *string to allow null values
}
