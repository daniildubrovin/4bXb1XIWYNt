package config

import "os"

type ServerConfig struct {
	Addr string
}

type BDConfig struct {
	Pass   string
	User   string
	Addr   string
	DBName string
}

type Config struct {
	Server ServerConfig
	BD     BDConfig
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Server: ServerConfig{
			Addr: getEnv("SERVER_ADDR", "localhost:8080"),
		},
		BD: BDConfig{
			Pass:   getEnv("POSTGRES_PASSWORD", ""),
			User:   getEnv("POSTGRES_USER", ""),
			Addr:   getEnv("POSTGRES_ADDR", "localhost:5432"),
			DBName: getEnv("POSTGRES_DB", "postgres"),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
