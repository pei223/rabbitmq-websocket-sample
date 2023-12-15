package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/pei223/rabbitmq-websocket-sample/internal/logger"
	"github.com/pei223/rabbitmq-websocket-sample/internal/modules/messages"
	"github.com/pei223/rabbitmq-websocket-sample/internal/queue"
	"github.com/pei223/rabbitmq-websocket-sample/internal/sessions"
	"golang.org/x/net/websocket"
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

	// モジュール生成
	sessionManager := sessions.New()
	messageSessionMod := messages.NewMessageSessionModule(ctx, sessionManager)
	worker := messages.NewMessageWorker(sessionManager, amqpClient)

	// ワーカーを別goroutineで実行
	go func() {
		logger.Info().Msg("Start Worker")
		err := worker.Run(ctx)
		if err != nil {
			logger.Warn().Err(err).Msg("Stop Worker")
			return
		}
		logger.Info().Msg("Stop Worker")
	}()

	mux := http.NewServeMux()

	// ルーティング定義
	// http.Handle("/ws", websocket.Handler(messageSessionMod.HandleMessageSession))
	// 本来上のコードが良いがクロスオリジンの関係で403になる. これだと通るらしい.
	// https://qiita.com/m0a/items/f6405bc29073a7609050
	mux.HandleFunc("/ws",
		func(w http.ResponseWriter, req *http.Request) {
			s := websocket.Server{Handler: websocket.Handler(messageSessionMod.HandleMessageSession)}
			s.ServeHTTP(w, req.WithContext(ctx))
		})

	// http.ListenAndServeだとシャットダウンができないのでserverオブジェクトを用いる.
	server := &http.Server{
		Addr:    ":8100",
		Handler: mux,
	}
	// goroutineで元のcontextの終了待ちをしてからシャットダウンする.
	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error().Err(err).Msg("Error on shutdown")
		}
	}()

	logger.Info().Msg("Start Websocket server")
	if err := server.ListenAndServe(); err != nil {
		logger.Error().Err(err).Msg("Stop server error")
		return
	}
	logger.Info().Msg("Stop server")
}
