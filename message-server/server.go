package messageserver

import (
	"encoding/json"
	"errors"

	kafka "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/robot-gitee-software-package/community"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
)

func Init(service app.MessageService, handler community.EventHandler) messageServer {
	return messageServer{
		service: service,
		handler: giteeEventHandler{
			handler: handler,
		},
	}
}

type messageServer struct {
	service app.MessageService
	handler giteeEventHandler
}

func (m *messageServer) Subscribe(cfg *Config) error {
	subscribers := map[string]kafka.Handler{
		cfg.Topics.NewPkg:         m.handleNewPkg,
		cfg.Topics.CommunityEvent: m.handler.handle,
	}

	return kafka.Subscribe(cfg.GroupName, subscribers)
}

func (m *messageServer) handleNewPkg(msg []byte) error {
	if len(msg) == 0 {
		return errors.New("unexpect message: The payload is empty")
	}

	var v msgToHandleNewPkg
	if err := json.Unmarshal(msg, &v); err != nil {
		return err
	}

	cmd := v.toCmd()
	return m.service.NewPkg(&cmd)
}
