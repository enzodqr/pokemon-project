package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Generic function that encodes any given value and returns it as a HTTP response
// TODO: Remove the any type and define a struct for the data
func Encode[T any](w http.ResponseWriter, status int, value T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(value);
	if err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

type Validator interface {
	Valid(ctx context.Context) map[string]string
}

// Decode a give value and validates if the response meets required structure
func Decode[T Validator](r *http.Request) (T, map[string]string, error) {
	var value T

	err := json.NewDecoder(r.Body).Decode(&value)
	if err != nil {
		return value, nil, fmt.Errorf("decode json: %w", err)
	}

	problems := value.Valid(r.Context())
	if len(problems) > 0 {
		return value, problems, fmt.Errorf("invalid %T: %d problems", value, len(problems))
	}

	return value, nil, nil
}

