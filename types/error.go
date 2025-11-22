package types

// ErrorResponse represents the generic error object returned by Nobitex.
// Different endpoints may include:
//   - status: always present on failure
//   - code: machine-readable error identifier
//   - message: human-readable explanation
//   - detail: optional additional context or validation details
//
// The struct is intentionally permissive because Nobitex does not use a single
// unified schema across all endpoints.
type ErrorResponse struct {
	// Status indicates failure state, typically "failed".
	Status string `json:"status"`

	// Code is an optional short identifier describing the error category.
	Code string `json:"code,omitempty"`

	// Message provides a human-readable description of the problem.
	Message string `json:"message,omitempty"`

	// Detail may contain additional explanation or field-specific validation errors.
	Detail string `json:"detail,omitempty"`
}
