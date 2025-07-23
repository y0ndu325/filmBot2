package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"film_botV2/internal/config"
	"film_botV2/internal/database"
	"film_botV2/internal/handlers"
	"film_botV2/internal/service"
	"film_botV2/internal/tmdb"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if cfg.BotToken == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	db, err := database.NewDatabase(cfg.DBConfig.DSN())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	clearConfig := tgbotapi.NewUpdate(-1)
	clearConfig.Timeout = 0
	if _, err := bot.GetUpdates(clearConfig); err != nil {
		log.Printf("Error clearing updates queue: %v", err)
	}
	tmdbClient, err := tmdb.NewClient(cfg.TMDBAPIKey)
	if err != nil {
		log.Fatalf("Failed to create TMDB client: %v", err)
	}

	movieService := service.New(db, tmdbClient)
	handler := handlers.New(bot, movieService, cfg)

	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 6
	updateConfig.AllowedUpdates = []string{"message", "callback_query"}

	updates := bot.GetUpdatesChan(updateConfig)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for update := range updates {
			log.Printf("Received update: %+v", update)

			if update.Message != nil {
				log.Printf("Processing message: %s", update.Message.Text)
				if err := handler.HandleMessage(&update); err != nil {
					log.Printf("Error handling message: %v", err)
					// Отправляем пользователю сообщение об ошибке
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при обработке сообщения. Попробуйте еще раз.")
					if _, sendErr := bot.Send(msg); sendErr != nil {
						log.Printf("Error sending error message: %v", sendErr)
					}
				}
			} else if update.CallbackQuery != nil {
				log.Printf("Processing callback query: %s", update.CallbackQuery.Data)
				if err := handler.HandleCallback(&update); err != nil {
					log.Printf("Error handling callback_query: %v", err)
					chatID := update.CallbackQuery.Message.Chat.ID
					msg := tgbotapi.NewMessage(chatID, "Произошла ошибка при обработке нажатия. Попробуйте еще раз.")
					if _, sendErr := bot.Send(msg); sendErr != nil {
						log.Printf("Error sending callback error message: %v", sendErr)
					}
				}
				ack := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
				if _, ackErr := bot.Request(ack); ackErr != nil {
					log.Printf("Error sending callback acknowledgement: %v", ackErr)
				}
			} else {
				log.Printf("Received unknown update type")
			}
		}
	}()

	<-sigChan
	log.Println("Shutting down...")
}
