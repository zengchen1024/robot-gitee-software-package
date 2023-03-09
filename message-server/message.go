package messageserver

import "github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"

type messageOfNewPkg struct{}

func (msg *messageOfNewPkg) toCmd() (app.CmdToCreatePR, error) {
	return app.CmdToCreatePR{}, nil
}
