package community

import (
	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
)

func (impl *eventHandler) handlePRState(e *sdk.PullRequestEvent) error {
	_, err := impl.repo.Find(int(e.Number))
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

		return impl.service.HandlePRMerged(&cmd)

	case sdk.StatusClosed:
		updateBy := e.GetUpdatedBy().GetLogin()
		if impl.cli.robotName == updateBy {
			return nil
		}

		cmd := impl.closedPrCmd(e.Number, updateBy)

		return impl.service.HandlePRClosed(&cmd)

	default:
		return nil
	}
}

func (impl *eventHandler) closedPrCmd(num int64, rejectedBy string) app.CmdToHandlePRClosed {
	return app.CmdToHandlePRClosed{
		PRNum:      int(num),
		RejectedBy: rejectedBy,
	}
}
