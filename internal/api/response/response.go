package response

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Internal = func(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusBadRequest, map[string]string{"error": "internal error"})
}

func JSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	render.Status(r, status)
	render.JSON(w, r, v)
}

func ValidationError(w http.ResponseWriter, r *http.Request, errs validator.ValidationErrors) {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	render.Status(r, http.StatusUnprocessableEntity)
	render.JSON(w, r, strings.Join(errMsgs, ", "))
}
