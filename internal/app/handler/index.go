package apphandler

import (
	"net/http"

	"github.com/gorilla/mux"
)

type appHandler struct {
	Router *mux.Router
	//Service *services.Service
}

func NewAppHandler(router *mux.Router) *appHandler {
	return &appHandler{
		Router: router,
		//Service: service,
	}
}

func (h *appHandler) RegisterHandlers() {
	h.Router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("spotify authorization succeded"))
	}).Methods("GET")
}
