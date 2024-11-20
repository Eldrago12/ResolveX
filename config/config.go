package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	RedisURL         string
	ServerPort       string
	RateLimit        int
	RateLimitTTL     int
	PrefetchInterval int
	PrefetchLimit    int
}

func LoadConfig() Config {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL not set")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	rateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil {
		rateLimit = 100
	}

	rateLimitTTL, err := strconv.Atoi(os.Getenv("RATE_LIMIT_TTL"))
	if err != nil {
		rateLimitTTL = 60
	}

	prefetchInterval, err := strconv.Atoi(os.Getenv("PREFETCH_INTERVAL"))
	if err != nil {
		prefetchInterval = 300
	}

	prefetchLimit, err := strconv.Atoi(os.Getenv("PREFETCH_LIMIT"))
	if err != nil {
		prefetchLimit = 10
	}

	return Config{
		RedisURL:         redisURL,
		ServerPort:       serverPort,
		RateLimit:        rateLimit,
		RateLimitTTL:     rateLimitTTL,
		PrefetchInterval: prefetchInterval,
		PrefetchLimit:    prefetchLimit,
	}
}
