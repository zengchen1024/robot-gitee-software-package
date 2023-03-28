package main

import (
	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

func (bot *robot) handlePRState(e *sdk.PullRequestEvent) error {
	_, err := bot.repo.Find(int(e.Number))
	if err != nil {
		if repository.IsErrorResourceNotFound(err) {
			err = nil
		}

		return err
	}

	switch e.GetPullRequest().GetState() {
	case sdk.StatusMerged:
		cmd := app.CmdToHandlePRMerged{
			PRNum: int(e.Number),
		}

		return bot.prService.HandlePRMerged(&cmd)
	case sdk.StatusClosed:
		r, err := bot.cli.GetBot()
		if err != nil {
			return err
		}

		updateBy := e.GetUpdatedBy().GetLogin()
		if r.Login == updateBy {
			return nil
		}

		cmd := bot.closedPrCmd(e.Number, updateBy)
		return bot.prService.HandlePRClosed(&cmd)
	default:
		return nil
	}
}

func (bot *robot) closedPrCmd(num int64, rejectedBy string) app.CmdToHandlePRClosed {
	return app.CmdToHandlePRClosed{
		PRNum:      int(num),
		RejectedBy: rejectedBy,
	}
}
