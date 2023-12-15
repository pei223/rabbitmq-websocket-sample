package messages

import (
	"net/http"

	"github.com/pei223/rabbitmq-websocket-sample/internal/logger"
	"github.com/pei223/rabbitmq-websocket-sample/internal/queue"
	"github.com/pei223/rabbitmq-websocket-sample/internal/views"
)

type MessageModule interface {
	PostMessage(w http.ResponseWriter, r *http.Request)
}

type messageModule struct {
	amqpClient *queue.AMQPClient
}

func NewMessageModule(amqpClient *queue.AMQPClient) MessageModule {
	return &messageModule{
		amqpClient: amqpClient,
	}
}

func (m *messageModule) PostMessage(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With().Logger()
	params, err := views.ParseBodyAndValidate[sampleMessageParams](r)
	if err != nil {
		views.RenderError(w, err)
		return
	}

	logger.Debug().Msg("publish message")

	queueSampleMessage := params.toQueueParam()
	err = m.amqpClient.PublishSampleMessage(&queueSampleMessage)
	if err != nil {
		logger.Warn().Err(err).Interface("message", queueSampleMessage).Msg("failed to publish")
		views.RenderError(w, views.NewUnexpectedErr(err.Error()))
		return
	}

	views.Render(w, http.StatusOK)
}
