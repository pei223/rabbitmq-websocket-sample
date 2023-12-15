package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/pei223/rabbitmq-websocket-sample/internal/logger"
	"github.com/pei223/rabbitmq-websocket-sample/internal/modules/messages"
	"github.com/pei223/rabbitmq-websocket-sample/internal/queue"
)

func main() {
	logger := logger.Logger.With().Logger()
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to read config")
		return
	}

	// キャンセル処理伝播
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger = logger.With().Interface("config", cfg).Logger()

	amqpClient, err := queue.NewClient(cfg.RabbitMqURL)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to setup amqp client")
		return
	}
	defer amqpClient.Close()

	messageMod := messages.NewMessageModule(amqpClient)

	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Route("/messages", func(r chi.Router) {
			r.Post("/", messageMod.PostMessage)
		})
	})
	server := &http.Server{
		Addr:    ":8090", // TODO 環境変数
		Handler: router,
	}
	// goroutineで元のcontextの終了待ちをしてからシャットダウンする.
	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error().Err(err).Msg("Error on shutdown")
		}
	}()
	logger.Info().Msg("Start API server")
	if err := server.ListenAndServe(); err != nil {
		logger.Error().Err(err).Msg("Stop server error")
		return
	}
	logger.Info().Msg("Stop server")
}
