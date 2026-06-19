package handlers

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"economy-go/database"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{Name: "balance", Description: "Баланс"},
	{Name: "deposit", Description: "Положить в банк", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionInteger, Name: "amount", Description: "Сумма", Required: true},
	}},
	{Name: "withdraw", Description: "Снять из банка", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionInteger, Name: "amount", Description: "Сумма", Required: true},
	}},
	{Name: "work", Description: "Заработать"},
	{Name: "pay", Description: "Перевести", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "Кому", Required: true},
		{Type: discordgo.ApplicationCommandOptionInteger, Name: "amount", Description: "Сумма", Required: true},
	}},
	{Name: "give_coins", Description: "Выдать монеты (админ)", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "Кому", Required: true},
		{Type: discordgo.ApplicationCommandOptionInteger, Name: "amount", Description: "Сумма", Required: true},
	}},
	{Name: "casino", Description: "Рискнуть", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionInteger, Name: "amount", Description: "Сумма", Required: true},
	}},
	{Name: "shop", Description: "Магазин цветов"},
	{Name: "buy_color", Description: "Купить цвет", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionString, Name: "color", Description: "Название цвета", Required: true},
	}},
	{Name: "paint", Description: "Применить купленный цвет"},
	{Name: "leaderboard", Description: "Топ 10"},
	{Name: "help", Description: "Помощь"},
}

func Register(s *discordgo.Session) {
	for _, cmd := range commands {
		if _, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd); err != nil {
			fmt.Println("Ошибка создания команды:", err)
		}
	}
}

func Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	userID := i.Member.User.ID

	switch data.Name {
	case "balance":
		balanceCmd(s, i, userID)
	case "deposit":
		depositCmd(s, i, userID, data)
	case "withdraw":
		withdrawCmd(s, i, userID, data)
	case "work":
		workCmd(s, i, userID)
	case "pay":
		payCmd(s, i, userID, data)
	case "give_coins":
		giveCmd(s, i, data)
	case "casino":
		casinoCmd(s, i, userID, data)
	case "shop":
		shopCmd(s, i)
	case "buy_color":
		buyCmd(s, i, data)
	case "paint":
		paintCmd(s, i, userID)
	case "leaderboard":
		leaderCmd(s, i)
	case "help":
		helpCmd(s, i)
	}
}

func balanceCmd(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	database.CreateUser(userID)
	user := database.GetUser(userID)

	embed := &discordgo.MessageEmbed{
		Title:       "💰 Баланс",
		Description: fmt.Sprintf("**Наличные:** %d монет\n**Банк:** %d монет", user.Balance, user.Bank),
		Color:       0xFFD700,
		Footer:      &discordgo.MessageEmbedFooter{Text: "/work — заработать"},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func depositCmd(s *discordgo.Session, i *discordgo.InteractionCreate, userID string, data discordgo.ApplicationCommandInteractionData) {
	var amount int
	for _, opt := range data.Options {
		if opt.Name == "amount" {
			amount = int(opt.IntValue())
			break
		}
	}

	if amount <= 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Сумма должна быть больше 0", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	database.CreateUser(userID)
	user := database.GetUser(userID)

	if amount > user.Balance {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Недостаточно монет!", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	database.UpdateBalance(userID, -amount)
	database.UpdateBank(userID, amount)

	embed := &discordgo.MessageEmbed{
		Title:       "🏦 Вклад",
		Description: fmt.Sprintf("Ты положил **%d** монет в банк", amount),
		Color:       0x2ECC71,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func withdrawCmd(s *discordgo.Session, i *discordgo.InteractionCreate, userID string, data discordgo.ApplicationCommandInteractionData) {
	var amount int
	for _, opt := range data.Options {
		if opt.Name == "amount" {
			amount = int(opt.IntValue())
			break
		}
	}

	if amount <= 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Сумма должна быть больше 0", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	database.CreateUser(userID)
	user := database.GetUser(userID)

	if amount > user.Bank {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Недостаточно монет в банке!", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	database.UpdateBalance(userID, amount)
	database.UpdateBank(userID, -amount)

	embed := &discordgo.MessageEmbed{
		Title:       "🏦 Снятие",
		Description: fmt.Sprintf("Ты снял **%d** монет из банка", amount),
		Color:       0x3498DB,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func workCmd(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	database.CreateUser(userID)
	user := database.GetUser(userID)

	if time.Since(user.LastWork).Hours() < 4 {
		left := 4 - time.Since(user.LastWork).Hours()
		embed := &discordgo.MessageEmbed{
			Title:       "⏳ Подожди",
			Description: fmt.Sprintf("Отдохни **%.1f** часов.", left),
			Color:       0xF1C40F,
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	salary := rand.Intn(150) + 50
	database.UpdateBalance(userID, salary)
	database.UpdateWork(userID)

	embed := &discordgo.MessageEmbed{
		Title:       "💼 Работа",
		Description: fmt.Sprintf("Ты заработал **%d** монет!", salary),
		Color:       0x2ECC71,
		Footer:      &discordgo.MessageEmbedFooter{Text: "Приходи через 4 часа"},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func payCmd(s *discordgo.Session, i *discordgo.InteractionCreate, userID string, data discordgo.ApplicationCommandInteractionData) {
	var targetID string
	var amount int
	for _, opt := range data.Options {
		if opt.Name == "user" {
			targetID = opt.UserValue(s).ID
		}
		if opt.Name == "amount" {
			amount = int(opt.IntValue())
		}
	}

	if amount <= 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Сумма > 0", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}
	if targetID == userID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Нельзя себе", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	database.CreateUser(userID)
	database.CreateUser(targetID)
	sender := database.GetUser(userID)
	if amount > sender.Balance {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Нет монет", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	database.UpdateBalance(userID, -amount)
	database.UpdateBalance(targetID, amount)

	embed := &discordgo.MessageEmbed{
		Title:       "💸 Перевод",
		Description: fmt.Sprintf("Переведено **%d** монет <@%s>", amount, targetID),
		Color:       0x3498DB,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func giveCmd(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	if !isAdmin(s, i) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Нет прав", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	var targetID string
	var amount int
	for _, opt := range data.Options {
		if opt.Name == "user" {
			targetID = opt.UserValue(s).ID
		}
		if opt.Name == "amount" {
			amount = int(opt.IntValue())
		}
	}

	database.CreateUser(targetID)
	database.UpdateBalance(targetID, amount)

	embed := &discordgo.MessageEmbed{
		Title:       "💸 Выдача",
		Description: fmt.Sprintf("<@%s> получил **%d** монет", targetID, amount),
		Color:       0x2ECC71,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func casinoCmd(s *discordgo.Session, i *discordgo.InteractionCreate, userID string, data discordgo.ApplicationCommandInteractionData) {
	var amount int
	for _, opt := range data.Options {
		if opt.Name == "amount" {
			amount = int(opt.IntValue())
		}
	}

	if amount <= 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Сумма > 0", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	database.CreateUser(userID)
	user := database.GetUser(userID)
	if amount > user.Balance {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Нет монет", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	var result string
	var color int
	if rand.Intn(2) == 0 {
		database.UpdateBalance(userID, amount)
		result = fmt.Sprintf("🎉 Выиграл **%d** монет!", amount)
		color = 0x2ECC71
	} else {
		database.UpdateBalance(userID, -amount)
		result = fmt.Sprintf("💀 Проиграл **%d** монет!", amount)
		color = 0xE74C3C
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎰 Казино",
		Description: result,
		Color:       color,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func shopCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	colors := []string{
		"🔴 Red", "🔵 Blue", "🟢 Green", "🟡 Yellow",
		"🟣 Purple", "🟠 Orange", "🩷 Pink", "🩵 Cyan",
		"⬜ White", "⬛ Black",
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎨 Магазин цветов",
		Description: "**Доступные цвета (5000 монет):**\n" + strings.Join(colors, "\n"),
		Color:       0x9B59B6,
		Footer:      &discordgo.MessageEmbedFooter{Text: "/buy_color [название]"},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func buyCmd(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	var colorName string
	for _, opt := range data.Options {
		if opt.Name == "color" {
			colorName = opt.StringValue()
		}
	}

	userID := i.Member.User.ID
	database.CreateUser(userID)
	user := database.GetUser(userID)

	if user.Balance < 5000 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Нужно 5000 монет", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	colorMap := map[string]int{
		"red": 0xE74C3C, "blue": 0x3498DB, "green": 0x2ECC71,
		"yellow": 0xF1C40F, "purple": 0x9B59B6, "orange": 0xE67E22,
		"pink": 0xE91E63, "cyan": 0x1ABC9C, "white": 0xFFFFFF, "black": 0x2C3E50,
	}

	colorHex := colorMap[strings.ToLower(colorName)]
	if colorHex == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Неверное название цвета", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	database.UpdateBalance(userID, -5000)
	database.SetColor(userID, colorName)

	roleName := "🎨 " + colorName
	role, err := s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
		Name:  roleName,
		Color: &colorHex,
	})
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Ошибка создания роли", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	s.GuildMemberRoleAdd(i.GuildID, userID, role.ID)

	embed := &discordgo.MessageEmbed{
		Title:       "✅ Куплен цвет",
		Description: fmt.Sprintf("Ты купил цвет **%s** и получил роль!", colorName),
		Color:       colorHex,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func paintCmd(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	user := database.GetUser(userID)
	if user.Color == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Нет купленного цвета. Купи в /shop", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	roleName := "🎨 " + user.Color
	roles, _ := s.GuildRoles(i.GuildID)
	for _, r := range roles {
		if r.Name == roleName {
			s.GuildMemberRoleAdd(i.GuildID, userID, r.ID)
			break
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎨 Цвет применён",
		Description: fmt.Sprintf("Ты применил цвет **%s**!", user.Color),
		Color:       0x9B59B6,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func leaderCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	users := database.GetAllUsers()

	sort.Slice(users, func(a, b int) bool {
		return users[a].Balance > users[b].Balance
	})

	var desc string
	limit := 10
	if len(users) < limit {
		limit = len(users)
	}

	for i := 0; i < limit; i++ {
		user := users[i]
		emoji := "👤"
		if i == 0 {
			emoji = "🏆"
		} else if i == 1 {
			emoji = "🥈"
		} else if i == 2 {
			emoji = "🥉"
		}
		desc += fmt.Sprintf("%s **%d.** <@%s> — **%d** монет\n", emoji, i+1, user.ID, user.Balance)
	}

	if desc == "" {
		desc = "Нет данных"
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🏆 Топ 10 по монетам",
		Description: desc,
		Color:       0xFFD700,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func helpCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title: "📋 Команды",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "💰 Экономика", Value: "/balance — баланс\n/deposit [сумма] — в банк\n/withdraw [сумма] — из банка\n/work — заработать\n/pay [user] [сумма] — перевести", Inline: false},
			{Name: "🎰 Развлечения", Value: "/casino [сумма] — рискнуть", Inline: false},
			{Name: "🎨 Цвета", Value: "/shop — магазин\n/buy_color [название] — купить\n/paint — применить", Inline: false},
			{Name: "👑 Админ", Value: "/give_coins [user] [сумма] — выдать", Inline: false},
			{Name: "🏆 Топ", Value: "/leaderboard — топ 10", Inline: false},
		},
		Color: 0x9B59B6,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Flags: discordgo.MessageFlagsEphemeral},
	})
}

func isAdmin(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	perms, _ := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	return perms&discordgo.PermissionAdministrator != 0
}
