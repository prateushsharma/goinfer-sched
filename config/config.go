package config

import (
	"os"
	"strconv"
)

// Config holds all runtime configuration for GoInferSched.
// Values are read from environment variables with sensible defaults.
type Config struct {
	Server    ServerConfig
	Scheduler SchedulerConfig
	Redis     RedisConfig
	Postgres  PostgresConfig
}

type ServerConfig struct {
	HTTPPort string // port the HTTP gateway listens on
	GRPCPort string // port for gRPC (Phase 2)
}

type SchedulerConfig struct {
	PlannerMode      string  // "heuristic" | "llm" | "hybrid"
	PlannerTimeoutMs int     // fall back to heuristic if planner takes longer than this
	AgingThresholdS  int     // seconds before a tier-3 request gets promoted to tier-2
	MaxBatchSize     int     // max requests grouped into one batch
	FlushDeadlineMs  int     // dispatch batch early if oldest request is this old (ms)
	VRAMSafetyMargin float64 // refuse a node if free VRAM < this fraction (0.15 = 15%)
	HealthIntervalMs int     // how often node agents report GPU stats (ms)
	MinRetryTokens   int     // requeue if fewer than this many tokens were streamed
}

type RedisConfig struct {
	Addr     string // e.g. "localhost:6379"
	Password string // empty string = no auth
}

type PostgresConfig struct {
	DSN string // full postgres connection string
}

// Load reads config from environment variables.
// Every field has a default so the app works out of the box locally.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			HTTPPort: getEnv("HTTP_PORT", "8080"),
			GRPCPort: getEnv("GRPC_PORT", "9090"),
		},
		Scheduler: SchedulerConfig{
			PlannerMode:      getEnv("PLANNER_MODE", "heuristic"),
			PlannerTimeoutMs: getEnvInt("PLANNER_TIMEOUT_MS", 50),
			AgingThresholdS:  getEnvInt("AGING_THRESHOLD_S", 30),
			MaxBatchSize:     getEnvInt("MAX_BATCH_SIZE", 8),
			FlushDeadlineMs:  getEnvInt("FLUSH_DEADLINE_MS", 200),
			VRAMSafetyMargin: getEnvFloat("VRAM_SAFETY_MARGIN", 0.15),
			HealthIntervalMs: getEnvInt("HEALTH_INTERVAL_MS", 500),
			MinRetryTokens:   getEnvInt("MIN_RETRY_TOKENS", 20),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Postgres: PostgresConfig{
			DSN: getEnv("POSTGRES_DSN", "postgres://user:pass@localhost/goinfer?sslmode=disable"),
		},
	}
}

// --- helpers ---

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvFloat(key string, defaultVal float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return defaultVal
}
