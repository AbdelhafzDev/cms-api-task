package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(LoadConfig),
)

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	var missing []string

	cfg := &Config{
		App: AppConfig{
			Name:         getEnv("APP_NAME", "cms-api"),
			Env:          getEnv("APP_ENV", "development"),
			Debug:        getEnvBool("APP_DEBUG", true),
			Version:      getEnv("APP_VERSION", "1.0.0"),
			AssetBaseURL: getEnv("ASSET_BASE_URL", ""),
		},
		HTTP: HTTPConfig{
			Host:            getEnv("HTTP_HOST", "0.0.0.0"),
			Port:            getEnvInt("HTTP_PORT", 8080),
			ReadTimeout:     getEnvDuration("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getEnvDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getEnvDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getEnvDuration("HTTP_SHUTDOWN_TIMEOUT", 30*time.Second),
			AllowedOrigins:  getEnvSlice("HTTP_ALLOWED_ORIGINS", []string{"*"}),
		},
		GRPC: GRPCConfig{
			Host:            getEnv("GRPC_HOST", "0.0.0.0"),
			Port:            getEnvInt("GRPC_PORT", 9090),
			MaxRecvMsgSize:  getEnvInt("GRPC_MAX_RECV_MSG_SIZE", 4194304),
			MaxSendMsgSize:  getEnvInt("GRPC_MAX_SEND_MSG_SIZE", 4194304),
			ShutdownTimeout: getEnvDuration("GRPC_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "postgres"),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "cms"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "cms"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			PublicKeyPath:      getEnv("AUTH_PUBLIC_KEY_PATH", ""),
			PrivateKeyPath:     getEnv("AUTH_PRIVATE_KEY_PATH", ""),
			AccessTokenExpiry:  getEnvDuration("AUTH_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getEnvDuration("AUTH_REFRESH_TOKEN_EXPIRY", 168*time.Hour),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			Output:     getEnv("LOG_OUTPUT", "stdout"),
			TimeFormat: getEnv("LOG_TIME_FORMAT", "2006-01-02T15:04:05.000Z07:00"),
		},
		Telemetry: TelemetryConfig{
			Enabled:      getEnvBool("TELEMETRY_ENABLED", false),
			OTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			SampleRate:   getEnvFloat("OTEL_SAMPLE_RATE", 1.0),
		},
		Search: SearchConfig{
			Host:      getEnv("MEILI_HOST", "meilisearch"),
			Port:      getEnvInt("MEILI_PORT", 7700),
			MasterKey: getEnv("MEILI_MASTER_KEY", ""),
		},
		Worker: WorkerConfig{
			PollInterval: getEnvDuration("WORKER_POLL_INTERVAL", 5*time.Second),
			BatchSize:    getEnvInt("WORKER_BATCH_SIZE", 10),
			MaxAttempts:  getEnvInt("WORKER_MAX_ATTEMPTS", 5),
		},
		Cache: CacheConfig{
			Host:     getEnv("REDIS_HOST", "redis"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
	}

	if cfg.IsProduction() {
		missing = append(missing, requireEnv("DB_PASSWORD")...)
		missing = append(missing, requireEnv("AUTH_PUBLIC_KEY_PATH")...)
		missing = append(missing, requireEnv("AUTH_PRIVATE_KEY_PATH")...)
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func requireEnv(keys ...string) []string {
	var missing []string
	for _, key := range keys {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	return missing
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func getEnvSlice(key string, fallback []string) []string {
	if v := os.Getenv(key); v != "" {
		return strings.Split(v, ",")
	}
	return fallback
}
