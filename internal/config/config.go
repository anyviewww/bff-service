package config

import "os"

type Config struct {
	MenuServiceAddr  string
	OrderServiceAddr string
	ServerPort       string
	JWTSecretKey     string
}

func Load() *Config {
	return &Config{
		MenuServiceAddr:  getEnv("MENU_SERVICE_ADDR", "menu-service:50051"),
		OrderServiceAddr: getEnv("ORDER_SERVICE_ADDR", "order-service:50052"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		JWTSecretKey:     getEnv("JWT_SECRET_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
