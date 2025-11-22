package utils

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// StructToURLParams converts a struct to a URL-encoded query string.
//
// This function uses the `json` struct tags as parameter keys and excludes
// fields with empty or zero values. It supports various data types, including
// slices, arrays, integers, floats, booleans, and strings.
//
// Supported Behavior:
//   - Fields with `json` tags are used as keys. Fields without tags or with
//     `json:"-"` are ignored.
//   - Zero values (e.g., empty strings, 0 for integers, 0.0 for floats) are omitted.
//   - Slices and arrays are converted to multiple key-value pairs.
//
// Parameters:
//   - inputStruct: The input struct to be converted into URL parameters. It
//     must be of kind `struct`.
//
// Returns:
//   - A URL-encoded query string as a `string`.
//   - An `error` if the input is not a struct or if any other issue occurs.
//
// Example:
//
//	type MyStruct struct {
//	    Name    string   `json:"name"`
//	    Age     int      `json:"age"`
//	    Tags    []string `json:"tags"`
//	    IsAdmin bool     `json:"is_admin"`
//	}
//
//	data := MyStruct{
//	    Name:    "John",
//	    Age:     30,
//	    Tags:    []string{"golang", "developer"},
//	    IsAdmin: true,
//	}
//
//	query, err := StructToURLParams(data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(query)
//	// Output: name=John&age=30&tags=golang&tags=developer&is_admin=true
//
// Limitations:
//   - Only fields with `json` tags are considered.
//   - Non-struct input will result in an error.
func StructToURLParams(inputStruct interface{}) (string, error) {
	values := url.Values{}

	// Get the type and value of the input struct
	v := reflect.ValueOf(inputStruct)
	t := reflect.TypeOf(inputStruct)

	// Ensure the input is a struct
	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("input must be a struct")
	}

	// Iterate through the struct fields
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Use the "json" tag if available; otherwise, use the field name
		key := field.Tag.Get("json")
		if key == "" || key == "-" {
			continue // Skip fields without a "json" tag or explicitly ignored
		}

		// Skip zero values
		if !value.IsValid() || value.IsZero() {
			continue
		}

		// Handle different kinds of fields
		switch value.Kind() {
		case reflect.Slice, reflect.Array:
			if value.Len() > 0 { // Only add non-empty slices/arrays
				for j := 0; j < value.Len(); j++ {
					values.Add(key, fmt.Sprintf("%v", value.Index(j).Interface()))
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() != 0 { // Skip zero values
				values.Add(key, strconv.FormatInt(value.Int(), 10))
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() != 0 { // Skip zero values
				values.Add(key, strconv.FormatFloat(value.Float(), 'f', -1, 64))
			}
		case reflect.Bool:
			values.Add(key, strconv.FormatBool(value.Bool())) // Always add booleans
		default:
			if value.String() != "" { // Skip empty strings
				values.Add(key, fmt.Sprintf("%v", value.Interface()))
			}
		}
	}

	// Encode and return the URL parameters
	return values.Encode(), nil
}
