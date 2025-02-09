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

type Bot struct {
	session    *discordgo.Session
	db         *storage.Database
	translator *translate.Translator
}

func NewBot(token string, db *storage.Database, translator *translate.Translator) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{session: s, db: db, translator: translator}
	s.AddHandler(bot.handleMessage)

	return bot, nil
}

func (b *Bot) Start() error {
	log.Println("Starting Discord bot...")
	return b.session.Open()
}

func (b *Bot) Stop() {
	log.Println("Stopping Discord bot...")
	b.session.Close()
}

func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	log.Printf("Received message: %s", m.Content)

	botID := s.State.User.ID
	content := strings.TrimSpace(strings.ReplaceAll(m.Content, "<@"+botID+">", ""))
	args := strings.Fields(content)

	if len(args) == 0 {
		return
	}

	command := args[0]

	switch command {
	case "!recent":
		lang, err := b.db.GetUserLanguage(m.Author.ID)
		if err != nil {
			lang = "en"
		}

		changes, err := stream.GetRecentChanges()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error fetching recent changes: "+err.Error())
			return
		}

		translatedChanges, err := b.translator.TranslateText(changes, lang)
		if err != nil {
			translatedChanges = changes
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Recent changes (%s): %s", lang, translatedChanges))

	case "!setLang":
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: !setLang [language_code]")
			return
		}
		lang := args[1]
		err := b.db.SetUserLanguage(m.Author.ID, lang)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to set language: "+err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Language set to "+lang)

	case "!help":
		lang, err := b.db.GetUserLanguage(m.Author.ID)
		if err != nil {
			lang = "en"
		}

		helpText := "Available commands:\n!recent - Get recent Wikipedia changes\n!setLang [language_code] - Set your preferred language\n!help - Show this help message"
		translatedHelp, err := b.translator.TranslateText(helpText, lang)
		if err != nil {
			translatedHelp = helpText
		}
		s.ChannelMessageSend(m.ChannelID, translatedHelp)

	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command. Type !help for a list of available commands.")
	}
}

func (b *Bot) SendMessage(channelID, message string) {
	_, err := b.session.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}
