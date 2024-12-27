package main

import (
	"fmt"
	"log"
	"net/http"
	"spoty/configs"
	"spoty/internal/app/handlers"
	"spoty/internal/spotify"

	"github.com/gorilla/mux"
)

func main() {
	cfg := configs.LoadConfigs()

	router := mux.NewRouter()

	{
		spotifyService := spotify.NewSpotifyService(cfg)
		spotifyHandler := spotify.NewSpotifyHandler(spotifyService, router)
		spotifyHandler.RegisterSpotifyHandler()
	}

	{
		appHandler := handlers.NewAppHandler(router)
		appHandler.RegisterHandlers()
	}

	fmt.Println("server is listening on port 8080")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), router); err != nil {
		log.Fatalf("error while starting server: \n%v", err)
	}
}
