package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"economy-go/config"
	"economy-go/database"
	"economy-go/handlers"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config.Load()
	database.Init()

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("❌ Ошибка:", err)
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		handlers.Handle(s, i)
	})

	err = dg.Open()
	if err != nil {
		log.Fatal("❌ Ошибка подключения:", err)
	}
	defer dg.Close()

	handlers.Register(dg)

	fmt.Println("✅ Бот запущен!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("👋 Бот остановлен")
}
