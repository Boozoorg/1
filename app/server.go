package app

import (
	"1/client"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Server struct {
	mux           *http.ServeMux
	clientService *client.Service
}

func NewServeMux(Mux *http.ServeMux, AccountService *client.Service) *Server {
	return &Server{
		mux:           Mux,
		clientService: AccountService,
	}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	s.mux.HandleFunc("/client.getByName", s.SaveAccount)
}

func (s *Server) SaveAccount(writer http.ResponseWriter, request *http.Request) {
	Name := request.URL.Query().Get("name")

	item, err := s.clientService.Registration(request.Context(), Name)
	if errors.Is(err, client.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}
