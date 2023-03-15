package main

import (
	"strings"

	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
)

const labelLgtmPrefix = "lgtm-"

func (bot *robot) handlePRState(e *sdk.PullRequestEvent) error {
	switch e.GetPullRequest().GetState() {
	case sdk.StatusMerged:
		users := bot.usersOfLgtmLabel(e)
		cmd := bot.mergedPrCmd(e.Number, users)

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
		Reason:     "pr is closed by maintainer",
		RejectedBy: rejectedBy,
	}
}

func (bot *robot) mergedPrCmd(num int64, users []string) app.CmdToHandlePRMerged {
	return app.CmdToHandlePRMerged{
		PRNum:      int(num),
		ApprovedBy: users,
	}
}

func (bot *robot) usersOfLgtmLabel(e *sdk.PullRequestEvent) []string {
	var users []string

	for _, label := range e.PullRequest.GetLabels() {
		if !strings.HasPrefix(label.Name, labelLgtmPrefix) {
			continue
		}

		user := strings.TrimPrefix(label.Name, labelLgtmPrefix)
		users = append(users, user)
	}

	return users
}
