package messages

import (
	"github.com/go-playground/validator/v10"
	"github.com/pei223/rabbitmq-websocket-sample/internal/queue"
	"github.com/pei223/rabbitmq-websocket-sample/internal/views"
)

type sampleMessageParams struct {
	Content string
	To      string
}

func (p sampleMessageParams) Validate() []views.InvalidParam {
	errs := views.AppValidator.Struct(p)
	if errs == nil {
		return nil
	}
	return views.ToInvalidParams(errs.(validator.ValidationErrors))
}

func (p *sampleMessageParams) toQueueParam() queue.SampleMessage {
	return queue.SampleMessage{
		Content: p.Content,
		To:      p.To,
	}
}
