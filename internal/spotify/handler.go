package spotify

import (
	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
	router  *mux.Router
}

func NewSpotifyHandler(service *Service, router *mux.Router) *Handler {
	return &Handler{
		service: service,
		router:  router,
	}
}

func (s *Handler) RegisterSpotifyHandler() {
	spotifyHandler := s.router.PathPrefix("/spotify").Subrouter()

	spotifyHandler.HandleFunc("/login", s.service.RedirectToLogin).Methods("GET")

	spotifyHandler.HandleFunc("/callback", s.service.Callback).Methods("GET")

	spotifyHandler.HandleFunc("/refresh", s.service.RefreshToken).Methods("GET")
}
