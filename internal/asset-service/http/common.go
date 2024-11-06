package http

import (
	"net/http"

	devxHttp "github.com/xBlaz3kx/DevX/http"
)

func badRequest(err error) (int, devxHttp.ErrorPayload) {
	return http.StatusBadRequest, devxHttp.ErrorPayload{
		Error:       "bad request",
		Code:        00001,
		Description: err.Error(),
	}
}
