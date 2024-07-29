package api

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/stanislavCasciuc/atom-fit-go/internal/utils"
	"log"
	"net/http"
)

type Server struct {
	addr string
	db   *sqlx.DB
}

func NewServer(addr string, db *sqlx.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/", writeHello).Methods("GET")

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}

func writeHello(w http.ResponseWriter, r *http.Request) {
	err := utils.WriteJSON(w, http.StatusOK, map[string]string{"hello": "world"})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
