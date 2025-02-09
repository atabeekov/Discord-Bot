package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wikibot/internal/translate"

	"wikibot/internal/discord"
	"wikibot/internal/storage"
	"wikibot/internal/stream"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := storage.InitDB("bot.db")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	if db != nil {
		defer db.Close() // Только если db успешно инициализирована
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Missing DISCORD_TOKEN environment variable")
	}

	translator, err := translate.NewTranslator()
	if err != nil {
		log.Fatalf("Failed to initialize translator: %v", err)
	}
	defer translator.Close() // Закрываем переводчик при завершении программы

	bot, err := discord.NewBot(token, db, translator)
	if err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}

	go bot.Start()
	defer bot.Stop()

	// Start the Wikipedia stream
	go stream.PollRecentChanges(func(changes string) {
		bot.SendMessage("912731302040600588", changes)
	})

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down bot.")
}
