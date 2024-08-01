package users

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/stanislavCasciuc/atom-fit-go/internal/api/response"
	"github.com/stanislavCasciuc/atom-fit-go/internal/config"
	"github.com/stanislavCasciuc/atom-fit-go/internal/lib/email"
	"github.com/stanislavCasciuc/atom-fit-go/internal/lib/jwt"
	"github.com/stanislavCasciuc/atom-fit-go/internal/lib/logger/sl"
	"github.com/stanislavCasciuc/atom-fit-go/internal/models"
	"github.com/stanislavCasciuc/atom-fit-go/internal/store/users"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log/slog"
	"net/http"
)

type Handler struct {
	store users.UserStore
	log   *slog.Logger
	cfg   config.Config
}

func NewHandler(store users.UserStore, log *slog.Logger) *Handler {
	return &Handler{store: store, log: log, cfg: config.Envs}
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	const op = "users.router.HandleLogin"
	requestId := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", requestId),
	)

	var payload models.LoginUserPayload

	err := render.DecodeJSON(r.Body, &payload)
	if errors.Is(err, io.EOF) {
		resp.JSON(w, r, http.StatusUnprocessableEntity, map[string]string{"error": "empty payload"})
		log.Error("request is empty")
		return
	}
	if err != nil {
		log.Error("failed to decode payload", sl.Err(err))
		resp.JSON(w, r, http.StatusUnprocessableEntity, map[string]string{"error": "failed to decode payload"})
		return
	}
	log.Info("request body successfully decoded", slog.Any("request_id", requestId))

	if err := validator.New().Struct(payload); err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("invalid request", sl.Err(err))
		resp.ValidationError(w, r, validateErr)
		return
	}

	log.With(slog.String("email", payload.Email))

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(users.UserNotFound, err) {
			resp.JSON(w, r, http.StatusBadRequest, map[string]string{"error": users.UserNotFound.Error()})
			log.Error("users not found")
			return
		}
		resp.Internal(w, r)
		log.Error("failed to get users", sl.Err(err))
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(payload.Password))
	if err != nil {
		resp.JSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid credentials"})
		log.Error("invalid credentials", sl.Err(err))
	}

	token, err := jwt.NewToken(*u, h.cfg.JwtCfg.Exp, h.cfg.JwtCfg.Secret)
	if err != nil {
		resp.Internal(w, r)
		log.Error("cannot to create token", sl.Err(err))
	}

	resp.JSON(w, r, http.StatusOK, map[string]string{"token": token})

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
		resp.JSON(w, r, http.StatusUnprocessableEntity, map[string]string{"error": "empty payload"})
		log.Error("request is empty")
		return
	}
	if err != nil {
		log.Error("failed to decode payload", sl.Err(err))
		resp.JSON(w, r, http.StatusUnprocessableEntity, map[string]string{"error": "failed to decode payload"})
		return
	}

	log.Info("request body successfully decoded", slog.Any("request_id", requestId))

	if err := validator.New().Struct(payload); err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("invalid request", sl.Err(err))
		resp.ValidationError(w, r, validateErr)
		return
	}

	log = log.With(slog.String("email", payload.Email))

	passHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error to hash password", sl.Err(err))
		resp.Internal(w, r)
	}

	id, err := h.store.CreateUser(payload, passHash)
	if err != nil {
		if errors.Is(users.UserAlreadyExist, err) {
			log.Error("users already exist")

			resp.JSON(w, r, http.StatusBadRequest, map[string]string{"error": "users already exist"})
			return
		}
		log.Error("fail to save users", sl.Err(err))
		resp.Internal(w, r)
		return
	}

	resp.JSON(w, r, http.StatusCreated, map[string]int{"user_id": id})

	err = email.SendVerifyUser(payload.Username, payload.Email, "dhfn10378efh23874hnfouh1u43he7f8")
	if err != nil {
		log.Error("error to send email", sl.Err(err))
	}
}
