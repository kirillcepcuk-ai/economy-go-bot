# 🪙 Economy Bot — Go

Discord economy bot with bank, shop, casino and color roles.

## 🔧 Tech Stack

| Technology | Description |
|------------|-------------|
| Go 1.21 | Programming language |
| discordgo | Discord API wrapper |
| JSON | File-based storage |

## 🚀 Commands

### 💰 Economy
- `/balance` — check your balance
- `/deposit [amount]` — deposit to bank
- `/withdraw [amount]` — withdraw from bank
- `/work` — earn coins (cooldown: 4 hours)
- `/pay [user] [amount]` — send coins to someone

### 👑 Admin
- `/give_coins [user] [amount]` — give coins to user

### 🎰 Casino
- `/casino [amount]` — gamble your coins (50% win chance)

### 🎨 Colors
- `/shop` — color shop
- `/buy_color [name]` — buy a color role (5000 coins)
- `/paint` — apply your bought color

### 🏆 Leaderboard
- `/leaderboard` — top 10 richest users

### 📋 Help
- `/help` — show all commands

## 🛠️ Installation

1. Install Go 1.21+
2. Create `.env` file with your token:
DISCORD_TOKEN=your_token_here
3. Run:
go mod tidy
go run .

## 📁 Project Structure

economy-go/
├── main.go
├── config/
│   └── config.go
├── database/
│   └── db.go
├── handlers/
│   └── commands.go
├── .env
├── .gitignore
└── README.md