package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	MongoURI string
	MongoDB  string

	RedisAddr string
	RedisPass string
	RedisDB   int

	JWTSecret          string
	JWTDuration        time.Duration
	JWTRefreshDuration time.Duration

	TelegramToken string
	OpenAIKey     string
	AIHost        string
	AIModel       string
}

func Load() *Config {
	_ = godotenv.Load()

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	jwtDuration, err := time.ParseDuration(os.Getenv("JWT_DURATION"))
	if err != nil {
		panic("invalid JWT_DURATION")
	}

	jwtRefreshDuration, err := time.ParseDuration(os.Getenv("JWT_REFRESH_DURATION"))
	if err != nil {
		panic("invalid JWT_REFRESH_DURATION")
	}

	return &Config{
		AppPort:  os.Getenv("APP_PORT"),
		MongoURI: os.Getenv("MONGO_URI"),
		MongoDB:  os.Getenv("MONGO_DB"),

		RedisAddr: os.Getenv("REDIS_ADDR"),
		RedisPass: os.Getenv("REDIS_PASSWORD"),
		RedisDB:   redisDB,

		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTDuration:        jwtDuration,
		JWTRefreshDuration: jwtRefreshDuration,

		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		OpenAIKey:     os.Getenv("OPENAI_API_KEY"),
		AIHost:        os.Getenv("AI_HOST"),
		AIModel:       os.Getenv("AI_MODEL"),
	}
}
