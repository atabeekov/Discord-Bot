package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"wikibot/internal/discord"
	"wikibot/internal/storage"
	"wikibot/internal/stream"
	"wikibot/internal/translate"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := storage.InitDB("bot.db")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Missing DISCORD_TOKEN environment variable")
	}

	translator := translate.NewTranslator()

	bot, err := discord.NewBot(token, db, translator)
	if err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}

	go bot.Start()
	defer bot.Stop()

	go stream.PollRecentChanges(func(changes string) {
		translated, err := translator.TranslateText(changes, "ru")
		if err != nil {
			log.Printf("Translation error: %v", err)
			bot.SendMessage("912731302040600588", changes)
		} else {
			bot.SendMessage("912731302040600588", translated)
		}
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down bot.")
}
