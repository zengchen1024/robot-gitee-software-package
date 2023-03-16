package repositoryimpl

import (
	"time"

	"github.com/google/uuid"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

const (
	mergedStatus   = 1
	unMergedStatus = 2
)

type SoftwarePkgPRDO struct {
	// must set "uuid" as the name of column
	PkgId     uuid.UUID `gorm:"column:uuid;type:uuid"`
	Link      string    `gorm:"column:link"`
	PkgName   string    `gorm:"column:pkg_name"`
	Num       int       `gorm:"column:num"`
	Merged    int       `gorm:"column:merge"`
	CreatedAt int64     `gorm:"column:created_at"`
	UpdatedAt int64     `gorm:"column:updated_at"`
}

func (s softwarePkgPR) toSoftwarePkgPRDO(p *domain.PullRequest, id uuid.UUID, do *SoftwarePkgPRDO) {
	*do = SoftwarePkgPRDO{
		PkgId:     id,
		Num:       p.Num,
		Link:      p.Link,
		PkgName:   p.Pkg.Name,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if p.IsMerged() {
		do.Merged = mergedStatus
	} else {
		do.Merged = unMergedStatus
	}
}

func (do *SoftwarePkgPRDO) toDomainPullRequest() (pr domain.PullRequest) {
	pr.Link = do.Link
	pr.Num = do.Num

	if do.Merged == mergedStatus {
		pr.SetMerged()
	}

	pr.Pkg.Name = do.PkgName
	pr.Pkg.Id = do.PkgId.String()

	return
}
