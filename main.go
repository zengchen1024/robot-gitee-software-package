package main

import (
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

	if err := kafka.Init(&cfg.MQ, log); err != nil {
		logrus.Fatalf("init kafka failed, err:%s", err.Error())
	}
	defer kafka.Exit()

	// app
	msApp := newSoftwarePkgApp(cfg)

	// message server
	ms := messageserver.Init(msApp)
	if err := ms.Subscribe(&cfg.MessageServer); err != nil {
		logrus.Fatalf("start side car failed, err:%s", err.Error())
	}

	// start
	r := newRobot(c)

	framework.Run(r, o.service)
}

func newSoftwarePkgApp(cfg *config.Config) app.MessageService {
	return nil
}
