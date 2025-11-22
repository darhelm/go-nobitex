package nobitex

import (
	"encoding/json"
	"fmt"

	t "github.com/darhelm/go-nobitex/types"
)

// GoNobitexError is the base error for all client-side and API-level failures.
type GoNobitexError struct {
	Message string // high-level description
	Err     error  // wrapped underlying error (optional)
}

func (e *GoNobitexError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *GoNobitexError) Unwrap() error { return e.Err }

// RequestError indicates a failure in preparing, sending, or reading a request.
type RequestError struct {
	GoNobitexError

	// Operation describes what failed, such as "creating request"
	// "sending request", or "parsing response".
	Operation string
}

// APIError represents an error returned directly by the Nobitex API.
// It is fully aligned with the documented error structure:
//
//	{
//	  "status":  "failed",
//	  "code":    "SomeError",
//	  "message": "Description"
//	}
type APIError struct {
	GoNobitexError

	Status     string // "failed"
	Code       string // error identifier
	StatusCode int    // HTTP status code (400/401/429/etc.)
	Detail     string // optional extra field occasionally returned
}

// parseErrorResponse parses the response body into an APIError,
// falling back gracefully if the body is not valid JSON.
func parseErrorResponse(statusCode int, respBody []byte) *APIError {
	var body t.ErrorResponse
	_ = json.Unmarshal(respBody, &body)

	var raw map[string]any
	_ = json.Unmarshal(respBody, &raw)

	var detail string
	if v, ok := raw["detail"].(string); ok {
		detail = v
	}

	// Determine the final error message to expose.
	msg := body.Message
	if msg == "" && detail != "" {
		msg = detail
	}
	if msg == "" {
		msg = fmt.Sprintf("API error (%d)", statusCode)
	}

	return &APIError{
		GoNobitexError: GoNobitexError{Message: msg},
		Status:         body.Status,
		Code:           body.Code,
		StatusCode:     statusCode,
		Detail:         detail,
	}
}
