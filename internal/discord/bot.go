package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"wikibot/internal/storage"
	"wikibot/internal/stream"
	"wikibot/internal/translate"
)

// Bot struct
type Bot struct {
	session    *discordgo.Session
	db         *storage.Database
	translator *translate.Translator
}

// NewBot initializes the Discord bot
func NewBot(token string, db *storage.Database, translator *translate.Translator) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{session: s, db: db, translator: translator}

	// Register message handler
	s.AddHandler(bot.handleMessage)

	return bot, nil
}

// Start runs the bot
func (b *Bot) Start() error {
	log.Println("Starting Discord bot...")
	return b.session.Open()
}

// Stop gracefully stops the bot
func (b *Bot) Stop() {
	log.Println("Stopping Discord bot...")
	b.session.Close()
}

// Handle incoming messages
func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return // Игнорируем сообщения от ботов
	}

	log.Printf("Received message: %s", m.Content)

	// Убираем упоминание бота из сообщения, если оно есть
	botID := s.State.User.ID
	content := strings.TrimSpace(strings.ReplaceAll(m.Content, "<@"+botID+">", ""))
	args := strings.Fields(content)

	if len(args) == 0 {
		log.Println("Empty message content.")
		return
	}

	command := args[0]
	log.Printf("Processing command: %s", command)

	switch command {
	case "!recent":
		lang, err := b.db.GetUserLanguage(m.Author.ID)
		if err != nil {
			log.Println("Error fetching user language:", err)
			lang = "en" // По умолчанию английский
		}

		// Получаем последние изменения
		changes, err := stream.GetRecentChanges()
		if err != nil {
			log.Println("Error fetching recent changes:", err)
			s.ChannelMessageSend(m.ChannelID, "Error fetching recent changes: "+err.Error())
			return
		}

		// Переводим изменения
		translatedChanges, err := b.translator.TranslateText(changes, lang)
		if err != nil {
			log.Println("Error translating changes:", err)
			translatedChanges = changes // Если перевод не удался, отправляем оригинал
		}

		// Отправляем переведённые изменения
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Recent changes (%s): %s", lang, translatedChanges))

	case "!setLang":
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: !setLang [language_code]")
			return
		}
		lang := args[1]

		// Устанавливаем язык
		err := b.db.SetUserLanguage(m.Author.ID, lang)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to set language: "+err.Error())
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Language set to "+lang)

	case "!help":
		lang, err := b.db.GetUserLanguage(m.Author.ID)
		if err != nil {
			log.Println("Error fetching user language:", err)
			lang = "en" // По умолчанию английский
		}

		helpText := "Available commands:\n!recent - Get recent Wikipedia changes\n!setLang [language_code] - Set your preferred language\n!help - Show this help message"
		translatedHelp, err := b.translator.TranslateText(helpText, lang)
		if err != nil {
			log.Println("Error translating help text:", err)
			translatedHelp = helpText
		}
		s.ChannelMessageSend(m.ChannelID, translatedHelp)

	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command. Type !help for a list of available commands.")
	}
}

// SendMessage sends a message to a Discord channel
func (b *Bot) SendMessage(channelID, message string) {
	_, err := b.session.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}
