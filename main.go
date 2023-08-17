package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/config"
	"github.com/opensourceways/robot-gitee-software-package/message-server"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/codeimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/emailimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/messageimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/postgresql"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/pullrequestimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/repositoryimpl"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/watch"
	"github.com/opensourceways/robot-gitee-software-package/utils"
)

type options struct {
	service liboptions.ServiceOptions
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.Parse(args)

	return o
}

func main() {
	logrusutil.ComponentInit("software-package")
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.Errorf("Invalid options, err:%s", err.Error())

		return
	}

	// cfg
	cfg, err := config.LoadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Errorf("load config failed, err:%s", err.Error())

		return
	}

	// postgresql
	if err = postgresql.Init(&cfg.Postgresql.DB); err != nil {
		logrus.Errorf("init db failed, err:%s", err.Error())

		return
	}

	// encryption
	if err = utils.InitEncryption(cfg.Encryption.EncryptionKey); err != nil {
		logrus.Errorf("init encryption failed, err:%s", err.Error())

		return
	}

	// kafka
	if err = kafka.Init(&cfg.Kafka, log); err != nil {
		logrus.Errorf("init kafka failed, err:%s", err.Error())

		return
	}

	defer kafka.Exit()

	// run
	run(cfg)
}

func run(cfg *config.Config) {
	pullRequest, err := pullrequestimpl.NewPullRequestImpl(&cfg.PullRequest)
	if err != nil {
		logrus.Errorf("init pull request failed, err:%s", err.Error())

		return
	}

	repo := repositoryimpl.NewSoftwarePkgPR(&cfg.Postgresql.Config)

	packageService := app.NewPackageService(
		repo,
		messageimpl.NewMessageImpl(cfg.MessageServer.Message),
		emailimpl.NewEmailService(cfg.Email),
		pullRequest,
		codeimpl.NewCodeImpl(cfg.Code),
	)

	// message server
	err = messageserver.Init(
		&cfg.MessageServer,
		app.NewMessageService(repo, pullRequest),
	)
	if err != nil {
		logrus.Errorf("init message server failed, err:%s", err.Error())

		return
	}

	// watch
	w := watch.NewWatchingImpl(cfg.Watch, repo, packageService)
	w.Start()
	defer w.Stop()

	// wait
	wait()
}

func wait() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			logrus.Info("receive done. exit normally")
			return

		case <-sig:
			logrus.Info("receive exit signal")
			called = true
			done()
			return
		}
	}(ctx)

	<-ctx.Done()
}
