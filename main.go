package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"telebotmusicos/app/bot"
	"telebotmusicos/app/locales"
	"telebotmusicos/app/player"
	"telebotmusicos/app/userbot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	BotToken     string
	ApiID        int
	ApiHash      string
	UserbotPhone string
}

func loadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Log: No .env file found, reading from environment.")
	}
	apiID, _ := strconv.Atoi(os.Getenv("TELEGRAM_API_ID"))
	return Config{
		BotToken:     os.Getenv("BOT_TOKEN"),
		ApiID:        apiID,
		ApiHash:      os.Getenv("TELEGRAM_API_HASH"),
		UserbotPhone: os.Getenv("USERBOT_PHONE_NUMBER"),
	}, nil
}

func main() {
	log.Println("Log: Bot startup sequence initiated.")
	if err := locales.Load(); err != nil {
		log.Fatalf("Fatal: Failed to load language files: %v", err)
	}
	log.Println("Log: Locales loaded successfully.")

	cfg, _ := loadConfig()
	if cfg.BotToken == "" || cfg.ApiID == 0 || cfg.ApiHash == "" {
		log.Fatal("Fatal: One or more required environment variables are not set.")
	}
	log.Println("Log: Configuration loaded successfully.")

	// Create a context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	botAPI, err := bot.Initialize(cfg.BotToken)
	if err != nil {
		log.Fatalf("Fatal: Bot initialization failed: %v", err)
	}

	userbotClient, err := userbot.Initialize(ctx, cfg.ApiID, cfg.ApiHash, cfg.UserbotPhone)
	if err != nil {
		log.Fatalf("Fatal: Userbot initialization failed: %v", err)
	}

	playerService := player.New(botAPI, userbotClient)
	handler := &bot.Handler{Player: playerService}

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine for Bot
	go func() {
		defer wg.Done()
		log.Println("Log: Bot listener starting...")
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := botAPI.GetUpdatesChan(u)

		for {
			select {
			case update := <-updates:
				handler.HandleUpdate(botAPI, update)
			case <-ctx.Done():
				log.Println("Log: Shutting down bot listener.")
				botAPI.StopReceivingUpdates()
				return
			}
		}
	}()

	// Goroutine for Userbot
	go func() {
		defer wg.Done()
		if err := userbotClient.Start(ctx); err != nil {
			if err != context.Canceled {
				log.Printf("Error: Userbot client stopped with error: %v", err)
			}
		}
		log.Println("Log: Userbot client shut down.")
	}()

	log.Println("Log: Bot is running. Press Ctrl+C to exit.")
	<-ctx.Done() // Wait for interrupt signal

	stop() // Stop listening for signals
	log.Println("Log: Shutting down all services...")
	wg.Wait() // Wait for all goroutines to finish
	log.Println("Log: Shutdown complete.")
}