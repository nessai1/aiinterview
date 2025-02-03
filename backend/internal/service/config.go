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
	IsDev   bool

	Secret         string
	InvitationCode string

	PSQLAddress string // мне лень делать декомпозицию
}

func FetchConfigFromEnv() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("cannot load .env file: %w", err)
	}

	address := os.Getenv("SERVICE_ADDR")
	if address == "" {
		return Config{}, fmt.Errorf("missing env variable SERVICE_ADDR")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		return Config{}, fmt.Errorf("missing env variable SECRET")
	}

	invitationCode := os.Getenv("INVITATION_CODE")

	proxyUrl := os.Getenv("PROXY_URL")
	proxyLogin := os.Getenv("PROXY_LOGIN")
	proxyPassword := os.Getenv("PROXY_PASSWORD")

	isDevEnv := os.Getenv("DEV")
	isDev := false
	if isDevEnv == "Y" {
		isDev = true
	}

	psqlAddr := os.Getenv("PSQL_ADDR")
	if psqlAddr == "" {
		return Config{}, fmt.Errorf("missing env variable PSQL_ADDR")
	}

	proxyConfig := ai.Config{ProxyURL: proxyUrl, ProxyLogin: proxyLogin, ProxyPassword: proxyPassword}
	config := Config{
		Address:        address,
		Secret:         secret,
		InvitationCode: invitationCode,
		OpenAI:         proxyConfig,
		IsDev:          isDev,
		PSQLAddress:    psqlAddr,
	}

	return config, nil
}
