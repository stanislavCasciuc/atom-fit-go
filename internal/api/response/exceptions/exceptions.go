package exceptions

import "net/http"

var (
	Internal      = http.StatusInternalServerError
	BadRequest    = http.StatusBadRequest
	Unprocessable = http.StatusUnprocessableEntity
	Created       = http.StatusCreated
	Ok            = http.StatusOK
)
