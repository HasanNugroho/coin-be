package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	MongoURI string
	MongoDB  string

	RedisAddr string
	RedisPass string
	RedisDB   int

	JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load()

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	return &Config{
		AppPort:   os.Getenv("APP_PORT"),
		MongoURI: os.Getenv("MONGO_URI"),
		MongoDB:  os.Getenv("MONGO_DB"),

		RedisAddr: os.Getenv("REDIS_ADDR"),
		RedisPass: os.Getenv("REDIS_PASSWORD"),
		RedisDB:   redisDB,

		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
