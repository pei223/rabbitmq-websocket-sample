package messages

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pei223/rabbitmq-websocket-sample/internal/logger"
	"github.com/pei223/rabbitmq-websocket-sample/internal/sessions"
	"golang.org/x/net/websocket"
)

type MessageSessionModule interface {
	HandleMessageSession(ws *websocket.Conn)
}

type messageSessionModule struct {
	ctx            context.Context
	sessionManager *sessions.SessionManager
}

func NewMessageSessionModule(ctx context.Context, sessionManager *sessions.SessionManager) MessageSessionModule {
	return &messageSessionModule{
		ctx:            ctx,
		sessionManager: sessionManager,
	}
}

func (m *messageSessionModule) HandleMessageSession(ws *websocket.Conn) {
	userName := ""
	logger := logger.Logger.With().Logger()

	defer func() {
		if userName != "" {
			m.sessionManager.Delete(userName)
		}
		logger.Info().Str("name", userName).Msg("connection close")
		ws.Close()
	}()

	// 初回のメッセージを送信
	err := websocket.Message.Send(ws, "Server: Hello, Client!")
	if err != nil {
		logger.Error().Err(err).Msg("failed first send")
		return
	}

	for {
		logger.Info().Msg("session start")

		// defaultがあるのでブロックしない
		select {
		case <-m.ctx.Done():
			logger.Info().Msg("Context done")
			return
		default:
		}

		time.Sleep(500 * time.Millisecond)

		logger.Debug().Msg("wait for message")

		// Client からのメッセージを読み込む
		msg := ""
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Debug().Msg("close session EOF")
				return
			}
			logger.Error().Err(err).Msg("failed to read message")
			continue
		}
		logger.Debug().Str("message", msg).Msg("receive message")

		if !strings.HasPrefix(msg, "[username]") || userName != "" {
			// ユーザー名以外も受け付けてそのままクライアントに返す
			err := websocket.Message.Send(ws, fmt.Sprintf("Received: %s", msg))
			if err != nil {
				logger.Error().Err(err).Msg("failed send")
			}
			continue
		}

		// 特定のprefixがあればユーザー名と見なす.
		// ユーザー名をキーにしてセッションをグローバル保存する.
		userName = strings.TrimPrefix(msg, "[username]")
		m.sessionManager.Add(userName, ws)

		logger = logger.With().Str("username", userName).Logger()

		// Client からのメッセージを元に返すメッセージを作成し送信する
		err := websocket.Message.Send(ws, fmt.Sprintf("Username accepted. name: %s", userName))
		if err != nil {
			logger.Error().Err(err).Msg("failed to send")
		}
	}
}
