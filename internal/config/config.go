package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SlackBotToken      string
	SlackSigningSecret string
	DBPath             string
	BotPort            string
	WebPort            string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf(".env ファイルの読み込みに失敗: %w", err)
	}

	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN が設定されていません")
	}

	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if signingSecret == "" {
		return nil, fmt.Errorf("SLACK_SIGNING_SECRET が設定されていません")
	}

	return &Config{
		SlackBotToken:      token,
		SlackSigningSecret: signingSecret,
		DBPath:             os.Getenv("DB_PATH"),
		BotPort:            getEnvWithDefault("BOT_PORT", "3000"),
		WebPort:            getEnvWithDefault("WEB_PORT", "8080"),
	}, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
