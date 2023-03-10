package app

import "github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"

type CmdToCreatePR = domain.SoftwarePkg

type CmdToHandleCI struct {
	PRNum        int
	FailedReason string
}

func (c *CmdToHandleCI) isSuccess() bool {
	return c.FailedReason == ""
}
