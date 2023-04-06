package messageserver

import (
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/community"
)

const (
	eventTypePush      = "push"
	msgHeaderUUID      = "X-GitHub-Delivery"
	msgHeaderUserAgent = "User-Agent"
	msgHeaderEventType = "X-GitHub-Event"
)

type giteeEventHandler struct {
	userAgent string
	handler   community.EventHandler
}

func (msg *giteeEventHandler) handle(payload []byte, header map[string]string) error {
	eventType, err := msg.parseRequest(header)
	if err != nil {
		return fmt.Errorf("invalid msg, err:%s", err.Error())
	}

	if eventType != eventTypePush {
		return errors.New("not pushed event")
	}

	e := new(sdk.PullRequestEvent)
	if err = json.Unmarshal(payload, e); err != nil {
		return err
	}

	return msg.handler.HandlePREvent(e)
}

func (msg *giteeEventHandler) parseRequest(header map[string]string) (
	eventType string, err error,
) {
	if header == nil {
		err = errors.New("no header")

		return
	}

	if header[msgHeaderUserAgent] != msg.userAgent {
		err = errors.New("unknown " + msgHeaderUserAgent)

		return
	}

	if eventType = header[msgHeaderEventType]; eventType == "" {
		err = errors.New("missing " + msgHeaderEventType)

		return
	}

	if header[msgHeaderUUID] == "" {
		err = errors.New("missing " + msgHeaderUUID)
	}

	return
}
