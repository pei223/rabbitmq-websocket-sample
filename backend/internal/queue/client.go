package queue

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/pei223/rabbitmq-websocket-sample/internal/logger"
	"github.com/streadway/amqp"
)

var sampleMessageExchangeName = "samplemessageexchange"
var sampleMessageQueueName = "samplemessagequeue"

type AMQPClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type DeferFunc func()

func NewClient(url string) (*AMQPClient, error) {
	c := AMQPClient{}

	amqpUrl := fmt.Sprintf("amqp://%s", url)
	conn, err := amqp.DialTLS(amqpUrl, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	c.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("failed to channel to rabbitmq: %w", err)
	}
	c.ch = ch

	// kindはいまいち分からない
	// 永続の場合はdurableはtrue, autoDeleteはfalseが良さそう
	// nowaitはサーバーからの確認を待たずに実行する. trueにする理由は特になさそう.
	// internalはよく分からないのでfalseで
	// argsは知らん
	// https://pkg.go.dev/github.com/streadway/amqp#Channel.ExchangeDeclare
	// TODO: amqp周りexchange nameを共通にしたい
	// TODO: 各引数の解釈
	if err := ch.ExchangeDeclare(sampleMessageExchangeName, "fanout", true, false, false, false, nil); err != nil {
		c.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Queue作成
	// exclusiveは何が良いか分からない
	q, err := ch.QueueDeclare(sampleMessageQueueName, true, false, false, false, nil)
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("failed to declare queue rabbitmq: %w", err)
	}
	// QueueにExchangeをBind
	// keyはよく分からない
	if err := ch.QueueBind(q.Name, "", sampleMessageExchangeName, false, nil); err != nil {
		c.Close()
		return nil, fmt.Errorf("failed to bind queue rabbitmq: %w", err)
	}

	return &c, nil
}

func (c *AMQPClient) ConsumeSampleMessage() (<-chan amqp.Delivery, error) {
	// auto-ackは無効にする
	// exclusiveは何が良いか分からない
	// consumerもどういう意味があるのか分からない
	return c.ch.Consume(sampleMessageQueueName, "samplemessageconsumer", false, true, false, false, nil)
}

func (c *AMQPClient) PublishSampleMessage(message *SampleMessage) error {
	bytes, err := json.Marshal(&message)
	if err != nil {
		return fmt.Errorf("failed to marshal message on publish: %w", err)
	}
	return c.ch.Publish(sampleMessageExchangeName, "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        bytes,
		// メッセージを永続化
		DeliveryMode: 2,
	})
}

func (c *AMQPClient) Close() {
	logger := logger.Logger.With().Logger()

	if c.ch != nil {
		if err := c.ch.Close(); err != nil {
			logger.Warn().Err(err).Msg("failed to close channel")
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			logger.Warn().Err(err).Msg("failed to close conn")
		}
	}
}
