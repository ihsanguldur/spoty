package configs

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	SpotifyClintId      string
	SpotifyClientSecret string
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
			SpotifyClintId:      getEnv("SPOTIFY_CLIENT_ID", "spotify-client-id"),
			SpotifyClientSecret: getEnv("SPOTIFY_CLIENT_SECRET", "spotify-client-secret"),
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
