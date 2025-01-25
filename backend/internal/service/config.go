package service

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nessai1/aiinterview/internal/ai"
	"os"
)

type Config struct {
	Address string
	OpenAI  ai.Config

	DBConnectAddr string // мне лень делать декомпозицию
}

func FetchConfigFromEnv() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("cannot load .env file: %w", err)
	}

	address := os.Getenv("SERVICE_ADDR")
	if address == "" {
		return Config{}, fmt.Errorf("missing env variable SERVICE_ADDR")
	}

	proxyUrl := os.Getenv("PROXY_URL")
	proxyLogin := os.Getenv("PROXY_LOGIN")
	proxyPassword := os.Getenv("PROXY_PASSWORD")

	psqlAddr := os.Getenv("PSQL_ADDR")
	if psqlAddr == "" {
		return Config{}, fmt.Errorf("missing env variable PSQL_ADDR")
	}

	proxyConfig := ai.Config{ProxyURL: proxyUrl, ProxyLogin: proxyLogin, ProxyPassword: proxyPassword}
	config := Config{Address: address, OpenAI: proxyConfig}

	return config, nil
}
