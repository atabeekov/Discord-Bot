Wikibot – Wikipedia Recent Changes Discord Bot

Wikibot is a Discord bot that fetches recent changes from Wikipedia and sends updates to a specified channel. It supports language filtering and on-demand translation of Wikipedia changes.

Features

- Live Wikipedia Recent Changes Feed: Fetches and displays recent Wikipedia edits.
- Multi-language Support: Filters Wikipedia changes based on language.
- On-Demand Translation: Uses Google Translate API to translate changes.
- Database Storage: Saves user preferences for language filtering.
- Discord Bot Integration: Listens to commands and sends Wikipedia updates.

Installation

Prerequisites
- Go (latest version recommended)
- SQLite
- A Discord Bot Token

Steps
- Clone the repository:
git clone https://github.com/yourusername/wikibot.git

cd wikibot

- Create a .env file in the root directory and add:
DISCORD_TOKEN=your_discord_bot_token

- Install dependencies:
go mod tidy
Run the bot:
go run main.go
Usage

Discord Commands
!recent – Fetches the latest Wikipedia changes.
!setLang <lang_code> – Sets the preferred language for Wikipedia changes.
!help – Shows the available commands.
Example
!setLang en
!recent
Project Structure

wikibot/
│── internal/
│   ├── discord/      # Handles Discord bot logic
│   ├── storage/      # Database operations
│   ├── stream/       # Fetches Wikipedia changes
│   ├── translate/    # Handles text translation
│── main.go           # Entry point
│── go.mod            # Go module file
│── .env              # Environment variables (ignored in Git)
Dependencies

DiscordGo – Discord bot library
GORM – ORM for SQLite
gtranslate – Translation library
Contributing

Feel free to submit issues and pull requests.
