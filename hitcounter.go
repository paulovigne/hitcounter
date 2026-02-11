package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	enableRedis bool
	redisClient *redis.Client
	localHits   atomic.Uint64
)

func main() {

	// ===== MODO HEALTHCHECK (para Docker) =====
	if len(os.Args) > 1 && os.Args[1] == "--healthcheck" {
		if err := runHealthcheck(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	// ===== CONFIG =====
	enableRedis = parseBoolEnv("ENABLE_REDIS", false)

	if enableRedis {
		initRedis()
		log.Println("Redis ENABLED: using Redis to store hits")
	} else {
		log.Println("Redis DISABLED: using in-memory counter")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", hitHandler)
	mux.HandleFunc("/healthz", healthHandler)

	port := getEnv("PORT", "8080")

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("Listening on :%s\n", port)
	log.Fatal(server.ListenAndServe())
}

func parseBoolEnv(key string, defaultVal bool) bool {
	env := os.Getenv(key)
	if env == "" {
		return defaultVal
	}
	parsed, err := strconv.ParseBool(env)
	if err != nil {
		return defaultVal
	}
	return parsed
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func initRedis() {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := os.Getenv("REDIS_PASSWORD")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
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
		hits = localHits.Add(1)
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hits: %d\n", hits)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := runHealthcheck(); err != nil {
		http.Error(w, "unhealthy", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func incrementRedis() (uint64, error) {
	val, err := redisClient.Incr(context.Background(), "hit_counter").Result()
	return uint64(val), err
}

func runHealthcheck() error {

	// Se Redis estiver habilitado, valida conexão
	if enableRedis {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			return err
		}
	}

	// Valida se porta HTTP está aberta (opcional mas interessante)
	port := getEnv("PORT", "8080")
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 2*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}
