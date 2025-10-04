package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server       ServerConfig
	MongoDB      MongoDBConfig
	Redis        RedisConfig
	JWT          JWTConfig
	ExternalAPIs ExternalAPIsConfig
	OCR          OCRConfig
	Cache        CacheConfig
	FileUpload   FileUploadConfig
}

type ServerConfig struct {
	Port   string
	AppEnv string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type ExternalAPIsConfig struct {
	RestCountriesAPIURL string
	ExchangeRateAPIURL  string
	ExchangeRateAPIKey  string
}

type OCRConfig struct {
	Service       string
	TesseractPath string
	TempDir       string
}

type CacheConfig struct {
	DefaultTTL          time.Duration
	ExpenseListTTL      time.Duration
	PendingApprovalsTTL time.Duration
	CurrencyRateTTL     time.Duration
	OCRResultTTL        time.Duration
}

type FileUploadConfig struct {
	MaxFileSize int64
	UploadDir   string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port:   getEnv("PORT", "8080"),
			AppEnv: getEnv("APP_ENV", "development"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "expensio"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			AccessTokenExpiry:  parseDuration(getEnv("JWT_ACCESS_TOKEN_EXPIRY", "15m")),
			RefreshTokenExpiry: parseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRY", "168h")), // 7 days
		},
		ExternalAPIs: ExternalAPIsConfig{
			RestCountriesAPIURL: getEnv("RESTCOUNTRIES_API_URL", "https://restcountries.com/v3.1"),
			ExchangeRateAPIURL:  getEnv("EXCHANGERATE_API_URL", "https://api.exchangerate-api.com/v4/latest"),
			ExchangeRateAPIKey:  getEnv("EXCHANGERATE_API_KEY", ""),
		},
		OCR: OCRConfig{
			Service:       getEnv("OCR_SERVICE", "tesseract"),
			TesseractPath: getEnv("TESSERACT_PATH", "/usr/bin/tesseract"),
			TempDir:       getEnv("OCR_TEMP_DIR", "./tmp/ocr"),
		},
		Cache: CacheConfig{
			DefaultTTL:          time.Duration(getEnvAsInt("CACHE_DEFAULT_TTL", 900)) * time.Second,
			ExpenseListTTL:      time.Duration(getEnvAsInt("CACHE_EXPENSE_LIST_TTL", 900)) * time.Second,
			PendingApprovalsTTL: time.Duration(getEnvAsInt("CACHE_PENDING_APPROVALS_TTL", 300)) * time.Second,
			CurrencyRateTTL:     time.Duration(getEnvAsInt("CACHE_CURRENCY_RATE_TTL", 3600)) * time.Second,
			OCRResultTTL:        time.Duration(getEnvAsInt("CACHE_OCR_RESULT_TTL", 86400)) * time.Second,
		},
		FileUpload: FileUploadConfig{
			MaxFileSize: int64(getEnvAsInt("MAX_FILE_SIZE", 10485760)), // 10MB default
			UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
		},
	}

	AppConfig = config
	return config
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves environment variable as integer or returns default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// parseDuration parses duration string (e.g., "15m", "7d")
func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("Error parsing duration %s, using default 15m: %v", s, err)
		return 15 * time.Minute
	}
	return duration
}

// IsProduction checks if app is running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.AppEnv == "production"
}

// GetRedisAddr returns Redis address in host:port format
func (c *Config) GetRedisAddr() string {
	return c.Redis.Host + ":" + c.Redis.Port
}
