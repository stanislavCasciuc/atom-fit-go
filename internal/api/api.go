package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	mwLogger "github.com/stanislavCasciuc/atom-fit-go/internal/api/midleware/logger"
	"github.com/stanislavCasciuc/atom-fit-go/internal/config"
	"github.com/stanislavCasciuc/atom-fit-go/internal/lib/logger/sl"
	"github.com/stanislavCasciuc/atom-fit-go/internal/services/users"
	users2 "github.com/stanislavCasciuc/atom-fit-go/internal/store/users"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	db  *sqlx.DB
	log *slog.Logger
	cfg config.Config
}

func NewServer(db *sqlx.DB, log *slog.Logger) *Server {
	return &Server{
		db:  db,
		log: log,
		cfg: config.Envs,
	}
}

func (s *Server) Run() error {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(s.log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	userStore := users2.NewStore(s.db)
	userHandlers := users.NewHandler(userStore, s.log)

	router.Post("/api/register", userHandlers.HandleRegister)
	router.Post("/api/login", userHandlers.HandleLogin)
	router.Post("/api/activate", userHandlers.ActivateUserHandler)

	s.log.Info("Listening on", slog.String("addr", s.cfg.HttpServer.Addr))
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         s.cfg.HttpServer.Addr,
		Handler:      router,
		ReadTimeout:  s.cfg.Timeout,
		WriteTimeout: s.cfg.Timeout,
		IdleTimeout:  s.cfg.IdleTimout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			s.log.Error("failed to start server")
		}
	}()

	s.log.Info("server started")

	<-done
	s.log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.log.Error("failed to stop server", sl.Err(err))

		return err
	}

	// TODO: close storage

	s.log.Info("server stopped")
	return nil
}
