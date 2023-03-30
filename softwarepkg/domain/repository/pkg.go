package repository

import "github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"

type SoftwarePkg interface {
	Add(pkg *domain.SoftwarePkg) error
	Save(*domain.SoftwarePkg) error
	Find(int) (domain.SoftwarePkg, error)
	FindAll() ([]domain.SoftwarePkg, error)
	Remove(int) error
}
