package messageserver

import (
	"errors"
	"strconv"
	"strings"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
)

type messageOfNewPkg struct{}

func (msg *messageOfNewPkg) toCmd() (app.CmdToCreatePR, error) {
	return app.CmdToCreatePR{}, nil
}

type messageOfApprovedPkg struct {
	PkgId      string `json:"pkg_id"`
	PkgName    string `json:"pkg_name"`
	RelevantPR string `json:"pr"`
}

func (msg *messageOfApprovedPkg) toCmd() (cmd app.CmdToMergePR, err error) {
	sp := strings.Split(strings.TrimSuffix(msg.RelevantPR, "/"), "/")
	if len(sp) == 0 {
		err = errors.New("relevant pr is empty")
		return
	}

	prNumInt, err := strconv.Atoi(sp[len(sp)-1])
	if err != nil {
		return
	}

	cmd.PRNum = prNumInt

	return
}
