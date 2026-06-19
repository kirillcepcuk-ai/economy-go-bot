package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Token string

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env не найден")
	}

	Token = os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		log.Fatal("❌ DISCORD_TOKEN не задан")
	}
}
