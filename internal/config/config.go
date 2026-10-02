package config

import "os"

type Config struct {
	ManagerURL       string
	EnrollmentID     string
	EnrollmentToken  string
	ConfigPath       string
	EndpointName     string
	Region           string
	Country          string
	City             string
	ListenAddr       string
	Version          string
}

func Load() Config {
	return Config{
		ManagerURL:      os.Getenv("STARGATE_MANAGER_URL"),
		EnrollmentID:    os.Getenv("STARGATE_ENROLLMENT_ID"),
		EnrollmentToken: os.Getenv("STARGATE_ENROLLMENT_TOKEN"),
		ConfigPath:      envOr("STARGATE_ENDPOINT_CONFIG", "/etc/stargate-endpoint/config.json"),
		EndpointName:    os.Getenv("STARGATE_ENDPOINT_NAME"),
		Region:          os.Getenv("STARGATE_REGION"),
		Country:         os.Getenv("STARGATE_COUNTRY"),
		City:            os.Getenv("STARGATE_CITY"),
		ListenAddr:      envOr("STARGATE_ENDPOINT_LISTEN", "127.0.0.1:8090"),
		Version:         envOr("STARGATE_ENDPOINT_VERSION", "0.1.0"),
	}
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
