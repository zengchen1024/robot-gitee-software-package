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

func (s softwarePkgPR) toSoftwarePkgPRDO(p *domain.PullRequest, id uuid.UUID, do *SoftwarePkgPRDO) error {
	email, err := toEmailDO(p.ImporterEmail)
	if err != nil {
		return err
	}

	*do = SoftwarePkgPRDO{
		PkgId:         id,
		Num:           p.Num,
		Status:        p.Status,
		Link:          p.Link,
		PkgName:       p.Pkg.Name,
		ImporterName:  p.ImporterName,
		ImporterEmail: email,
		SpecURL:       p.SrcCode.SpecURL,
		SrcRPMURL:     p.SrcCode.SrcRPMURL,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	return nil
}

func (do *SoftwarePkgPRDO) toDomainPullRequest() (pr domain.PullRequest, err error) {
	if pr.ImporterEmail, err = toEmail(do.ImporterEmail); err != nil {
		return
	}

	pr.Link = do.Link
	pr.Num = do.Num
	pr.Status = do.Status
	pr.Pkg.Name = do.PkgName
	pr.Pkg.Id = do.PkgId.String()
	pr.ImporterName = do.ImporterName
	pr.SrcCode.SpecURL = do.SpecURL
	pr.SrcCode.SrcRPMURL = do.SrcRPMURL

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
