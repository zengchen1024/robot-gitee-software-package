package messageserver

import "github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/messageimpl"

type Config struct {
	GroupName string             `json:"group_name"    required:"true"`
	Topics    Topics             `json:"topics"`
	Message   messageimpl.Config `json:"message"`
}

type Topics struct {
	NewPkg string `json:"new_pkg" required:"true"`
}
