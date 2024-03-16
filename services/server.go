package services

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
	users []User
}


func newServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
		users: []User{},
	}

	return s
}

func (s *Server) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i User
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		i.ID = uuid.New()
		i.Followers = []string{}
		i.Following = []string{}
		i.Username = "user"
		s.users = append(s.users, i)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}