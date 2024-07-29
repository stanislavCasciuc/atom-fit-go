package users

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stanislavCasciuc/atom-fit-go/internal/types"
	"github.com/stanislavCasciuc/atom-fit-go/internal/utils"
	"log/slog"
	"net/http"
)

type Handler struct {
	store types.UserStore
	log   *slog.Logger
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// handle login
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	const op = "users.routers.handleRegister"

	log := h.log.With(slog.String("op", op))

	var payload types.RegisterUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, errors.New("invalid payload"))
		h.log.Warn("unprocessable entity")
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		h.log.Warn("invalid payload")
		utils.WriteError(w, http.StatusUnprocessableEntity, fmt.Errorf("invalid payload: %w", err))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil || u == nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("user already exist"))
		h.log.Warn("user already exist ")
		return
	}

	log = log.With(slog.String("email", u.Email))

	id, err := h.store.CreateUser(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		h.log.Error("internal erorr", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}
