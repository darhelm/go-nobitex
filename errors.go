package nobitex

import (
	"encoding/json"
	"fmt"

	t "github.com/darhelm/go-nobitex/types"
)

type GoNobitexError struct {
	Message string
	Err     error
}

func (e *GoNobitexError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *GoNobitexError) Unwrap() error { return e.Err }

// RequestError Error returned when a request cannot be created/sent/read.
type RequestError struct {
	GoNobitexError
	Operation string
}

// APIError represents *any* server-side error returned by Nobitex.
//
// Nobitex usually returns one of:
//
//   { "status": "failed", "code": "ErrCode", "message": "msg" }
//   { "detail": "something went wrong" }
//   { "status": "failed", "message": "msg", "detail": "extra" }
//   { ... totally undocumented garbage ... }
type APIError struct {
	GoNobitexError

	Status     string
	Code       string
	Message    string
	Detail     string
	StatusCode int

	// Map of all parsed key->values for inspection (similar to go-bitpin)
	Fields map[string][]string
}

// parseErrorResponse creates the most complete APIError possible.
// It attempts all documented + undocumented patterns.
func parseErrorResponse(statusCode int, respBody []byte) *APIError {
	apiErr := &APIError{
		StatusCode: statusCode,
		Fields:     make(map[string][]string),
	}

	// #1 — Attempt to parse official Nobitex error format
	var base t.ErrorResponse
	_ = json.Unmarshal(respBody, &base)

	if base.Status != "" {
		apiErr.Status = base.Status
	}
	if base.Code != "" {
		apiErr.Code = base.Code
		apiErr.Fields["code"] = []string{base.Code}
	}
	if base.Message != "" {
		apiErr.Message = base.Message
		apiErr.Fields["message"] = []string{base.Message}
	}

	// #2 — Parse raw JSON object for extra fields, including "detail"
	raw := map[string]any{}
	_ = json.Unmarshal(respBody, &raw)

	for k, v := range raw {
		switch val := v.(type) {
		case string:
			apiErr.Fields[k] = []string{val}
			if k == "detail" {
				apiErr.Detail = val
				if apiErr.Message == "" {
					apiErr.Message = val
				}
			}
		case []any:
			// convert []any -> []string
			strs := make([]string, 0, len(val))
			for _, item := range val {
				strs = append(strs, fmt.Sprintf("%v", item))
			}
			apiErr.Fields[k] = strs
		default:
			apiErr.Fields[k] = []string{fmt.Sprintf("%v", v)}
		}
	}

	// #3 — If message is still empty, fallback
	if apiErr.Message == "" {
		apiErr.Message = fmt.Sprintf("API error (%d)", statusCode)
	}

	apiErr.GoNobitexError.Message = apiErr.Message
	return apiErr
}
