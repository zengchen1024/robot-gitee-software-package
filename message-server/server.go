package messageserver

import (
	"encoding/json"
	"errors"

	"github.com/opensourceways/robot-gitee-software-package/kafka"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
)

func Init(service app.MessageService) messageServer {
	return messageServer{
		service: service,
	}
}

type messageServer struct {
	service app.MessageService
}

func (m *messageServer) Subscribe(cfg *Config) error {
	subscribers := map[string]kafka.Handler{
		cfg.Topics.NewPkg: m.handleNewPkg,
	}

	return kafka.Instance().Subscribe(cfg.GroupName, subscribers)
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
