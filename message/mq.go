package message

import (
	"github.com/opensourceways/kafka-lib/kafka"
	"github.com/opensourceways/kafka-lib/mq"
	"github.com/sirupsen/logrus"
)

func mqInit(address string) error {
	err := kafka.Init(
		mq.Addresses(address),
		mq.Log(logrus.WithField("module", "kfk")),
	)
	if err != nil {
		return err
	}

	return kafka.Connect()
}

func mqExit() {
	if err := kafka.Disconnect(); err != nil {
		logrus.Errorf("exit kafka, err:%v", err)
	}
}
