package http

import (
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func badRequest(err error) (int, ErrorResponse) {
	return http.StatusBadRequest, ErrorResponse{Error: "Bad request", Message: err.Error()}
}
