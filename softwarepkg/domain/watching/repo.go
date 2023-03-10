package watching

import "github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"

type Watching interface {
	Apply(*domain.SoftwarePkgRepo) error
}
