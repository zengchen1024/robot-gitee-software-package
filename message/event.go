package message

import (
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/kafka-lib/kafka"
	"github.com/opensourceways/kafka-lib/mq"
	"github.com/sirupsen/logrus"
)

var robotLogin string

type Event struct {
	group       string
	cli         iClient
	cfg         *config
	log         *logrus.Entry
	subscribers map[string]mq.Subscriber
}

type iClient interface {
	GetBot() (sdk.User, error)
	CreatePullRequest(org, repo, title, body, head, base string, canModify bool) (sdk.PullRequest, error)
}

func InitEvent(cfgFile, group string, cli iClient) (*Event, error) {
	cfg, err := loadConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	if err = mqInit(cfg.KafkaAddress); err != nil {
		return nil, err
	}

	e := &Event{
		cli:   cli,
		cfg:   cfg,
		group: group,
	}

	if err = e.subscribe(); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Event) Exit() {
	e.unsubscribe()

	mqExit()
}

func (e *Event) unsubscribe() {
	for k, v := range e.subscribers {
		if err := v.Unsubscribe(); err != nil {
			logrus.Errorf("failed to unsubscribe for topic:%s, err:%v", k, err)
		}
	}
}

func (e *Event) subscribe() error {
	e.subscribers = make(map[string]mq.Subscriber)

	s, err := kafka.Subscribe(e.cfg.Topics.NewPkg, e.group, e.newPkgHandle)
	if err != nil {
		return err
	}

	e.subscribers[s.Topic()] = s

	return nil
}

func (e *Event) validateMessage(msg *mq.Message) error {
	if msg == nil {
		return errors.New("get a nil msg from broker")
	}

	if len(msg.Body) == 0 {
		return errors.New("unexpect message: The payload is empty")
	}

	return nil
}

func (e *Event) newPkgHandle(event mq.Event) error {
	if err := e.validateMessage(event.Message()); err != nil {
		return err
	}

	e.createPR(event.Message())

	return nil
}

func (e *Event) createPR(msg *mq.Message) {
	var c CreatePRParam
	if err := json.Unmarshal(msg.Body, &c); err != nil {
		e.log.WithError(err).Error("unmarshal")
		return
	}

	e.log = logrus.WithFields(
		logrus.Fields{
			"msg": c,
		},
	)

	if err := c.initRepo(e.cfg); err != nil {
		e.log.WithError(err).Error("init repo")
		return
	}

	if err := c.newBranch(e.cfg); err != nil {
		e.log.WithError(err).Error("new branch")
		return
	}

	if err := c.modifyFiles(e.cfg); err != nil {
		e.log.WithError(err).Error("modify files")
		return
	}

	if err := c.commit(e.cfg); err != nil {
		e.log.WithError(err).Error("commit")
		return
	}

	if err := e.createPRWithApi(c); err != nil {
		e.log.WithError(err).Error("create with api")
		return
	}
}

func (e *Event) createPRWithApi(p CreatePRParam) error {
	robotName, err := e.getRobotLogin()
	if err != nil {
		return err
	}

	head := fmt.Sprintf("%s:%s", robotName, branchName(e.cfg.PR.BranchName, p.PkgName))
	pr, err := e.cli.CreatePullRequest(
		e.cfg.PR.Org, e.cfg.PR.Repo, prName(e.cfg.PR.PRName, p.PkgName),
		p.ReasonToImportPkg, head, "master", true,
	)
	if err != nil {
		return err
	}

	e.log.Infof("pr number is %d", pr.Number)

	return nil
}

func (e *Event) getRobotLogin() (string, error) {
	if robotLogin == "" {
		v, err := e.cli.GetBot()
		if err != nil {
			return "", err
		}

		robotLogin = v.Login
	}

	return robotLogin, nil
}

func branchName(branchName, pkgName string) string {
	return fmt.Sprintf(branchName, pkgName)
}

func prName(prName, pkgName string) string {
	return fmt.Sprintf(prName, pkgName)
}
