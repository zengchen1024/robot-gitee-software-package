package repositoryimpl

import (
	"github.com/google/uuid"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain/repository"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/postgresql"
)

type softwarePkgPR struct {
	cli dbClient
}

func NewSoftwarePkgPR(cfg *Config) repository.PullRequest {
	return softwarePkgPR{cli: postgresql.NewDBTable(cfg.Table.SoftwarePkgPR)}
}

func (s softwarePkgPR) Add(p *domain.PullRequest) error {
	u, err := uuid.Parse(p.Pkg.Id)
	if err != nil {
		return err
	}

	var do SoftwarePkgPRDO
	if err = s.toSoftwarePkgPRDO(p, u, &do); err != nil {
		return err
	}

	filter := SoftwarePkgPRDO{PkgId: u}

	return s.cli.Insert(&filter, &do)
}

func (s softwarePkgPR) Save(p *domain.PullRequest) error {
	u, err := uuid.Parse(p.Pkg.Id)
	if err != nil {
		return err
	}
	filter := SoftwarePkgPRDO{PkgId: u}

	var do SoftwarePkgPRDO
	if err = s.toSoftwarePkgPRDO(p, u, &do); err != nil {
		return err
	}

	return s.cli.UpdateRecord(&filter, &do)
}

func (s softwarePkgPR) Find(num int) (domain.PullRequest, error) {
	filter := SoftwarePkgPRDO{Num: num}

	var res SoftwarePkgPRDO
	if err := s.cli.GetRecord(&filter, &res); err != nil {
		if s.cli.IsRowNotFound(err) {
			err = repository.NewErrorResourceNotFound(err)
		}

		return domain.PullRequest{}, err
	}

	return res.toDomainPullRequest()
}

func (s softwarePkgPR) FindAll(isMerged bool) ([]domain.PullRequest, error) {
	filter := SoftwarePkgPRDO{}
	if isMerged {
		filter.Merged = mergedStatus
	} else {
		filter.Merged = unMergedStatus
	}

	var res []SoftwarePkgPRDO

	if err := s.cli.GetRecords(
		&filter,
		&res,
		postgresql.Pagination{},
		nil,
	); err != nil {
		return nil, err
	}

	var p = make([]domain.PullRequest, len(res))

	for i := range res {
		v, err := res[i].toDomainPullRequest()
		if err != nil {
			return nil, err
		}

		p[i] = v
	}

	return p, nil
}

func (s softwarePkgPR) Remove(num int) error {
	filter := SoftwarePkgPRDO{Num: num}

	return s.cli.DeleteRecord(&filter)
}
