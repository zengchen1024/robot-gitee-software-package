package app

import "github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"

type CmdToNewPkg = domain.SoftwarePkg

type CmdToHandleCI struct {
	PRNum        int
	FailedReason string
}

func (c *CmdToHandleCI) isSuccess() bool {
	return c.FailedReason == ""
}

type CmdToMergePR struct {
	PRNum int
}

type CmdToClosePR struct {
	PRNum  int
	Reason string
}

type CmdToHandlePRMerged struct {
	PRNum      int
	ApprovedBy []string
}

type CmdToHandlePRClosed struct {
	PRNum      int
	Reason     string
	RejectedBy string
}
