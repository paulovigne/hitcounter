package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx          = context.Background()
	enableRedis  = false
	redisClient *redis.Client
	localHits    uint64 = 0
)

func main() {
	// Read env ENABLE_REDIS (default: false)
	env := os.Getenv("ENABLE_REDIS")
	if env != "" {
		parsed, err := strconv.ParseBool(env)
		if err == nil {
			enableRedis = parsed
		}
	}

	if enableRedis {
		initRedis()
		log.Println("Redis ENABLED: using Redis to store hits")
	} else {
		log.Println("Redis DISABLED: using in-memory counter")
	}

	http.HandleFunc("/", hitHandler)
	http.HandleFunc("/healthz", healthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initRedis() {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	password := os.Getenv("REDIS_PASSWORD")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctxTimeout).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func hitHandler(w http.ResponseWriter, r *http.Request) {
	var hits uint64
	var err error

	if enableRedis {
		hits, err = incrementRedis()
		if err != nil {
			http.Error(w, "Redis error", http.StatusInternalServerError)
			return
		}
	} else {
		localHits++
		hits = localHits
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hits: %d\n", hits)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func incrementRedis() (uint64, error) {
	val, err := redisClient.Incr(ctx, "hit_counter").Result()
	return uint64(val), err
}
