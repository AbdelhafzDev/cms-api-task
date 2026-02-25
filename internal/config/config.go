package config

import (
	"fmt"
	"strconv"
	"time"
)

type Config struct {
	App       AppConfig
	HTTP      HTTPConfig
	GRPC      GRPCConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	Log       LogConfig
	Telemetry TelemetryConfig
	Search    SearchConfig
	Worker    WorkerConfig
	Cache     CacheConfig
}

type AppConfig struct {
	Name         string
	Env          string
	Debug        bool
	Version      string
	AssetBaseURL string
}

type HTTPConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	AllowedOrigins  []string
}

type GRPCConfig struct {
	Host            string
	Port            int
	MaxRecvMsgSize  int
	MaxSendMsgSize  int
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func (c DatabaseConfig) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

type JWTConfig struct {
	PublicKeyPath       string
	PrivateKeyPath      string
	AccessTokenExpiry   time.Duration
	RefreshTokenExpiry  time.Duration
}

type LogConfig struct {
	Level      string
	Format     string
	Output     string
	TimeFormat string
}

type TelemetryConfig struct {
	Enabled      bool
	OTLPEndpoint string
	SampleRate   float64
}

type SearchConfig struct {
	Host      string
	Port      int
	MasterKey string
}

func (c SearchConfig) Addr() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

type WorkerConfig struct {
	PollInterval time.Duration
	BatchSize    int
	MaxAttempts  int
}

type CacheConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func (c CacheConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development" || c.App.Env == "dev"
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production" || c.App.Env == "prod"
}
