package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"spoty/configs"
	apphandler "spoty/internal/app/handler"
	"spoty/internal/spotify"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	cfg := configs.LoadConfigs()

	// TODO: change gorilla/mux to chi
	router := mux.NewRouter()

	// TODO: add static file

	spotifyService := spotify.NewSpotifyService(cfg)
	spotifyHandler := spotify.NewSpotifyHandler(spotifyService, router)
	spotifyHandler.RegisterSpotifyHandler()

	appHandler := apphandler.NewAppHandler(router)
	appHandler.RegisterHandlers()

	port := cfg.AppPort
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Printf("server started at port %s\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error while starting server: \n%v", err)
			cancel()
		}

	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()

		select {
		case sig := <-signalChan:
			fmt.Printf("received signal: %s", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	<-ctx.Done()

	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownContext); err != nil {
		fmt.Printf("error while shutting down server: \n%v", err)
	} else {
		fmt.Printf("server shutdown gracefully")
	}

	wg.Wait()
}
