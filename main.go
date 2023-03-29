package main

import (
	"context"
	"errors"
	"flag"
	"os"

	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/opensourceways/robot-gitee-lib/framework"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/opensourceways/server-common-lib/secret"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/config"
	"github.com/opensourceways/robot-gitee-software-package/kafka"
	"github.com/opensourceways/robot-gitee-software-package/message-server"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/codeimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/emailimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/messageimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/postgresql"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/pullrequestimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/repositoryimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/watchingimpl"
	"github.com/opensourceways/robot-gitee-software-package/utils"
)

type options struct {
	service       liboptions.ServiceOptions
	gitee         liboptions.GiteeOptions
	MsgConfigFile string
}

func (o *options) Validate() error {
	if err := o.service.Validate(); err != nil {
		return err
	}

	if o.MsgConfigFile == "" {
		return errors.New("missing message config file")
	}

	return o.gitee.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.gitee.AddFlags(fs)
	o.service.AddFlags(fs)
	fs.StringVar(&o.MsgConfigFile, "msg-config-file", "", "Path to message config file.")

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit(botName)
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
	}

	secretAgent := new(secret.Agent)
	if err := secretAgent.Start([]string{o.gitee.TokenPath}); err != nil {
		logrus.WithError(err).Fatal("Error starting secret agent.")
	}

	defer secretAgent.Stop()

	c := client.NewClient(secretAgent.GetTokenGenerator(o.gitee.TokenPath))

	// side car
	cfg, err := config.LoadConfig(o.MsgConfigFile)
	if err != nil {
		logrus.Fatalf("load config failed, err:%s", err.Error())
	}

	if err = postgresql.Init(&cfg.Postgresql.DB); err != nil {
		logrus.Fatalf("init db failed, err:%s", err.Error())
	}

	if err = kafka.Init(&cfg.MQ, log); err != nil {
		logrus.Fatalf("init kafka failed, err:%s", err.Error())
	}
	defer kafka.Exit()

	if err = utils.InitEncryption(cfg.Encryption.EncryptionKey); err != nil {
		logrus.Errorf("init encryption failed, err:%s", err.Error())

		return
	}

	pullRequest, err := pullrequestimpl.NewPullRequestImpl(c, cfg.PullRequest)
	if err != nil {
		logrus.Errorf("init pullRequest failed, err:%s", err.Error())

		return
	}

	email := emailimpl.NewEmailService(cfg.Email)
	message := messageimpl.NewMessageImpl(cfg.MessageServer.Message)
	repo := repositoryimpl.NewSoftwarePkgPR(&cfg.Postgresql.Config)
	code := codeimpl.NewCodeImpl(cfg.Code)

	prService := app.NewPullRequestService(repo, message, email, pullRequest, code)
	messageService := app.NewMessageService(repo, pullRequest)

	watch := watchingimpl.NewWatchingImpl(cfg.Watch, c, repo, prService)
	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan struct{})
	go watch.Start(ctx, stop)
	defer func() {
		cancel()
		<-stop
		logrus.Info("watch exit normally")
	}()

	// message server
	ms := messageserver.Init(messageService)
	if err := ms.Subscribe(&cfg.MessageServer); err != nil {
		logrus.Fatalf("start side car failed, err:%s", err.Error())
	}

	// start
	r := newRobot(c, prService, repo, cfg.Watch.Org)

	framework.Run(r, o.service)
}
