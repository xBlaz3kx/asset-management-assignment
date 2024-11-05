package errors

import (
	"net/http"

	"github.com/xBlaz3kx/DevX/errors"
)

var (
	ErrAssetNotFound      = errors.New(1001, http.StatusNotFound, "Asset not found")
	ErrAssetAlreadyExists = errors.New(1002, http.StatusConflict, "Asset already exists")
	ErrValidation         = errors.New(1003, http.StatusBadRequest, "Validation error")
	ErrTimeRangeViolation = errors.New(1004, http.StatusBadRequest, "Invalid time range provided")
)
