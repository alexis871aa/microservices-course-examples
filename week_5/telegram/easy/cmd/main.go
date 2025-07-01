package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	// Токен можно получить у @BotFather в Telegram
	token := "TELEGRAM_BOT_TOKEN"

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем бота с опциями
	b, err := bot.New(token,
		// Устанавливаем обработчик по умолчанию для всех сообщений
		bot.WithDefaultHandler(defaultHandler),
	)
	if err != nil {
		log.Printf("Failed to create bot: %v\n", err)
		return
	}

	// Регистрируем обработчик команды /start
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)

	// Регистрируем обработчик для сообщения "ping" (точное совпадение)
	b.RegisterHandler(bot.HandlerTypeMessageText, "ping", bot.MatchTypeExact, pingHandler)

	// Регистрируем обработчик для сообщений, содержащих "hello" (частичное совпадение)
	b.RegisterHandler(bot.HandlerTypeMessageText, "hello", bot.MatchTypeContains, helloHandler)

	log.Println("Bot started successfully! Send /start to begin")

	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping bot...")
		cancel()
	}()

	// Запускаем бота
	b.Start(ctx)
	log.Println("Bot stopped")
}

// startHandler обрабатывает команду /start
func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Отправляем приветственное сообщение
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `🤖 Привет! Я простой ping-pong бот!

Доступные команды:
• Напиши "ping" - получишь "pong"
• Напиши что-то с "hello" - получишь приветствие
• Любое другое сообщение - получишь эхо

Попробуй написать "ping"!`,
	})
	if err != nil {
		log.Printf("Failed to send start message: %v\n", err)
	}
}

// pingHandler обрабатывает сообщение "ping"
func pingHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Отвечаем "pong" на "ping"
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "🏓 pong",
		// Отвечаем на исходное сообщение
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
	})
	if err != nil {
		log.Printf("Failed to send pong message: %v\n", err)
	}
}

// helloHandler обрабатывает сообщения, содержащие "hello"
func helloHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Получаем имя пользователя
	userName := update.Message.From.FirstName
	if userName == "" {
		userName = "друг"
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "👋 Привет, " + userName + "! Как дела?",
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
	})
	if err != nil {
		log.Printf("Failed to send hello message: %v\n", err)
	}
}

// defaultHandler обрабатывает все остальные сообщения (эхо-бот)
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Проверяем, что это текстовое сообщение
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	// Отправляем эхо сообщения
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "📢 Эхо: " + update.Message.Text,
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
	})
	if err != nil {
		log.Printf("Failed to send echo message: %v\n", err)
	}
}
