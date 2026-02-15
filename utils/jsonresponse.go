package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type APIResponse struct {
	Success bool           `json:"success"`
	Data    any            `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
	//Code int `json:"code,omitempty"`
}

func RespondWithError(w http.ResponseWriter, code int, message string, err error) {
	errResp := ErrorResponse{
		Code:    code,
		Message: message,
	}
	if err != nil {
		errResp.Error = err.Error()
	}

	response := APIResponse{
		Success: false,
		Error:   &errResp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func RespondWithJSON(w http.ResponseWriter, code int, data any) {
	response := APIResponse{
		Success: true,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
