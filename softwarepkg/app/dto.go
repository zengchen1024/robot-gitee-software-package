package app

import "github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"

type CmdToHandleNewPkg = domain.SoftwarePkg

type CmdToHandleCI struct {
	PRNum        int
	RepoLink     string
	FailedReason string
}

func (c *CmdToHandleCI) isSuccess() bool {
	return c.FailedReason == ""
}

func (c *CmdToHandleCI) isPkgExisted() bool {
	return c.RepoLink != ""
}

type CmdToHandlePRMerged struct {
	PRNum int
}

type CmdToHandlePRClosed struct {
	PRNum      int
	RejectedBy string
}
