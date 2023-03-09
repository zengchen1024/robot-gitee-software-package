package watching

import "github.com/opensourceways/robot-gitee-software-package/pullrequest/domain"

type Watching interface {
	Apply(*domain.SoftwarePkgRepo) error
}
