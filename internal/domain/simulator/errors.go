package simulator

import (
	"net/http"

	"github.com/xBlaz3kx/DevX/errors"
)

var (
	ErrNoConfigForAsset = errors.New(2001, http.StatusNotFound, "Config for asset not found")
	ErrConfigNotFound   = errors.New(2002, http.StatusNotFound, "Config not found")
	ErrConfigValidation = errors.New(2003, http.StatusBadRequest, "Validation error")
)
