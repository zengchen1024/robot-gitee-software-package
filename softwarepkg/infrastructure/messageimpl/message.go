package messageimpl

import (
	"github.com/opensourceways/robot-gitee-software-package/kafka"
	messageserver "github.com/opensourceways/robot-gitee-software-package/message-server"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/message"
)

func NewMessageImpl(topics messageserver.Topics) *MessageImpl {
	return &MessageImpl{
		topics: topics,
	}
}

type MessageImpl struct {
	topics messageserver.Topics
}

func (m *MessageImpl) NotifyCIResult(e message.EventMessage) error {
	return send(m.topics.CIPassed, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return kafka.Instance().Publish(topic, body)
}
