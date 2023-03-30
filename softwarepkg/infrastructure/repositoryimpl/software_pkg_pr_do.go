package repositoryimpl

import (
	"time"

	"github.com/google/uuid"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
	"github.com/opensourceways/robot-gitee-software-package/utils"
)

type SoftwarePkgPRDO struct {
	// must set "uuid" as the name of column
	PkgId         uuid.UUID `gorm:"column:uuid;type:uuid"`
	Link          string    `gorm:"column:link"`
	PkgName       string    `gorm:"column:pkg_name"`
	Num           int       `gorm:"column:num"`
	Status        string    `gorm:"column:status"`
	ImporterName  string    `gorm:"column:importer_name"`
	ImporterEmail string    `gorm:"column:importer_email"`
	SpecURL       string    `gorm:"column:spec_url"`
	SrcRPMURL     string    `gorm:"column:src_rpm_url"`
	CreatedAt     int64     `gorm:"column:created_at"`
	UpdatedAt     int64     `gorm:"column:updated_at"`
}

func (s softwarePkgPR) toSoftwarePkgPRDO(p *domain.SoftwarePkg, id uuid.UUID, do *SoftwarePkgPRDO) error {
	email, err := toEmailDO(p.Importer.Email)
	if err != nil {
		return err
	}

	*do = SoftwarePkgPRDO{
		PkgId:         id,
		Num:           p.PullRequest.Num,
		Status:        p.Status,
		Link:          p.PullRequest.Link,
		PkgName:       p.Name,
		ImporterName:  p.Importer.Name,
		ImporterEmail: email,
		SpecURL:       p.Application.SourceCode.SpecURL,
		SrcRPMURL:     p.Application.SourceCode.SrcRPMURL,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	return nil
}

func (do *SoftwarePkgPRDO) toDomainPullRequest() (pkg domain.SoftwarePkg, err error) {
	if pkg.Importer.Email, err = toEmail(do.ImporterEmail); err != nil {
		return
	}

	pkg.PullRequest.Link = do.Link
	pkg.PullRequest.Num = do.Num
	pkg.Status = do.Status
	pkg.Name = do.PkgName
	pkg.Id = do.PkgId.String()
	pkg.Importer.Name = do.ImporterName
	pkg.Application.SourceCode.SpecURL = do.SpecURL
	pkg.Application.SourceCode.SrcRPMURL = do.SrcRPMURL

	return
}

func toEmailDO(email string) (string, error) {
	return utils.Encryption.Encrypt([]byte(email))
}

func toEmail(e string) (string, error) {
	v, err := utils.Encryption.Decrypt(e)
	if err != nil {
		return "", err
	}

	return string(v), nil
}
