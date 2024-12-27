package configs

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	SpotifyClientId     string
	SpotifyClientSecret string
	SpotifyRedirectUri  string
}

var once sync.Once
var instance *Config

func LoadConfigs() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Fatalf(".env file not found: \n%v", err)
		}

		instance = &Config{
			AppPort:             getEnv("APP_PORT", "8080"),
			SpotifyClientId:     getEnv("SPOTIFY_CLIENT_ID", "spotify-client-id"),
			SpotifyClientSecret: getEnv("SPOTIFY_CLIENT_SECRET", "spotify-client-secret"),
			SpotifyRedirectUri:  getEnv("SPOTIFY_REDIRECT_URI", "http://localhost:8080/spotify/callback"),
		}
	})

	return instance
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
