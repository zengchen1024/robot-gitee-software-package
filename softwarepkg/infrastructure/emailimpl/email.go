package emailimpl

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func NewEmailService(cfg Config) *emailImpl {
	return &emailImpl{
		cfg: cfg,
	}
}

type emailImpl struct {
	cfg Config
}

func (impl *emailImpl) Send(url string) error {
	d := gomail.NewDialer(
		impl.cfg.EmailServer.Host,
		impl.cfg.EmailServer.Port,
		impl.cfg.EmailServer.From,
		impl.cfg.EmailServer.AuthCode,
	)

	subject := "the CI of PR in openeuler/community is failed"
	content := fmt.Sprintf("the pr url: %s", url)

	message := gomail.NewMessage()
	message.SetHeader("From", impl.cfg.EmailServer.From)
	message.SetHeader("To", impl.cfg.MaintainerEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", content)

	if err := d.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
