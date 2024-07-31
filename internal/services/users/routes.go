package users

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/stanislavCasciuc/atom-fit-go/internal/api/response"
	"github.com/stanislavCasciuc/atom-fit-go/internal/lib/logger/sl"
	"github.com/stanislavCasciuc/atom-fit-go/internal/services/users/models"
	"io"
	"log/slog"
	"net/http"
)

type Handler struct {
	store UserStore
	log   *slog.Logger
}

func NewHandler(store UserStore, log *slog.Logger) *Handler {
	return &Handler{store: store, log: log}
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// handle login
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	const op = "users.routers.handleRegister"

	requestId := middleware.GetReqID(r.Context())

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", requestId),
	)

	var payload models.RegisterUserPayload

	err := render.DecodeJSON(r.Body, &payload)
	if errors.Is(err, io.EOF) {
		render.JSON(w, r, response.Error("empty payload"))
		log.Error("request is empty")
		return
	}
	if err != nil {
		log.Error("failed to decode payload", sl.Err(err))
		render.JSON(w, r, response.Error("failed to decode payload"))
		return
	}

	log.Info("request body successfully decoded", slog.Any("request_id", requestId))

	if err := validator.New().Struct(payload); err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, response.ValidationError(validateErr))
		return
	}

	log = log.With(slog.String("email", payload.Email))

	_, err = h.store.CreateUser(payload)
	if err != nil {
		if errors.Is(UserAlreadyExist, err) {
			log.Error("user already exist")
			render.JSON(w, r, response.Error("user already exist"))
			return
		}
		log.Error("fail to save user", sl.Err(err))
		render.JSON(w, r, response.Error("fail to save user"))
		return
	}

	render.JSON(w, r, response.OK())
}
