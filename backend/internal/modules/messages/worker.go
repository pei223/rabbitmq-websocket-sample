package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pei223/rabbitmq-websocket-sample/internal/logger"
	"github.com/pei223/rabbitmq-websocket-sample/internal/queue"
	"github.com/pei223/rabbitmq-websocket-sample/internal/sessions"
	"github.com/streadway/amqp"
)

type MessageWorker struct {
	sessionManager *sessions.SessionManager
	amqpClient     *queue.AMQPClient
}

func NewMessageWorker(sessionManager *sessions.SessionManager, amqpClient *queue.AMQPClient) *MessageWorker {
	return &MessageWorker{
		sessionManager: sessionManager,
		amqpClient:     amqpClient,
	}
}

func (w *MessageWorker) Run(ctx context.Context) error {
	logger := logger.Logger.With().Logger()

	// 受信開始
	msgs, err := w.amqpClient.ConsumeSampleMessage()
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	logger.Info().Msg("start consume")

	// RabbitMQからメッセージ取得したら、Toのユーザー名にメッセージを送信
CONSUMER_FOR:
	for {
		time.Sleep(500 * time.Millisecond)

		select {
		// defaultがないのでchannel受診するまでselect全体がブロックされる
		case <-ctx.Done():
			logger.Info().Msg("context done")
			break CONSUMER_FOR
		case msg, ok := <-msgs:
			if !ok {
				logger.Warn().Interface("msg", msg).Msg("Not ok")
				continue
			}
			err := w.onMessage(msg)
			if err != nil {
				logger.Warn().Err(err).Interface("msgBody", msg.Body).Msg("onMessage error")
				continue
			}
		}
	}
	return nil
}

func (w *MessageWorker) onMessage(msg amqp.Delivery) error {
	var message queue.SampleMessage
	if err := json.Unmarshal(msg.Body, &message); err != nil {
		// 形式がおかしいので破棄
		msg.Nack(false, false)
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	if err := w.sessionManager.Send(message.To, message.Content); err != nil {
		// 送信エラーはrequeue
		msg.Nack(false, true)
		return fmt.Errorf("failed to send message, %v: %w", message, err)
	}
	return nil
}
