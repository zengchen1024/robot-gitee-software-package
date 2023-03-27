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
