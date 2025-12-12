package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser             string
	DBPassword         string
	DBName             string
	DBHost             string
	DBPort             string
	JWTSecret          string
	RefreshTokenSecret string
	PushURL            string
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		panic("Impossible de charger le fichier .env")
	}

	return &Config{
		DBUser:             getEnv("DB_USER", "root"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "attendify"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "3312"),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", ""),
		PushURL:            getEnv("PUSH_URL", ""),
	}
}
