package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func SendAPIResponse(w http.ResponseWriter, status int, data any, cookies []http.Cookie) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			http.SetCookie(w, &cookie)
		}
	}
	return json.NewEncoder(w).Encode(generateAPIResponseBody(status, data))
}

func SendAPIErrorResponse(w http.ResponseWriter, status int, err error) {
	SendAPIResponse(w, status, err.Error(), nil)
}

func generateAPIResponseBody(status int, data any) map[string]any {
	return map[string]any{"status": status, "data": data}
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(v)
}
