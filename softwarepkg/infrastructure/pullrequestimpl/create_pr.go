package pullrequestimpl

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

func (impl *pullRequestImpl) createBranch(pkg *domain.SoftwarePkg) error {
	sigInfoFile := fmt.Sprintf(
		"sig/%s/sig-info.yaml",
		pkg.Application.ImportingPkgSig,
	)

	sigInfoData, err := impl.genAppendSigInfoData(pkg)
	if err != nil {
		return err
	}

	newRepoFile := fmt.Sprintf(
		"sig/%s/src-openeuler/%s/%s.yaml",
		pkg.Application.ImportingPkgSig,
		strings.ToLower(pkg.Name[:1]),
		pkg.Name,
	)

	newRepoData, err := impl.genNewRepoData(pkg)
	if err != nil {
		return err
	}

	params := []string{
		impl.cfg.ShellScript,
		impl.cfg.Robot.Username,
		impl.cfg.Robot.Token,
		impl.cfg.Robot.Email,
		impl.branchName(pkg.Name),
		impl.cfg.PR.Org,
		impl.cfg.PR.Repo,
		sigInfoFile,
		sigInfoData,
		newRepoFile,
		newRepoData,
	}

	out, err, _ := utils.RunCmd(params...)
	if err != nil {
		logrus.Errorf(
			"run create pr shell, err=%s, out=%s, params=%v",
			err.Error(), string(out), params[:len(params)-1],
		)
	}

	return err
}

func (impl *pullRequestImpl) branchName(pkgName string) string {
	return fmt.Sprintf("software_package_%s", pkgName)
}

func (impl *pullRequestImpl) genAppendSigInfoData(pkg *domain.SoftwarePkg) (string, error) {
	data := struct {
		PkgName       string
		ImporterEmail string
		Importer      string
	}{
		PkgName:       pkg.Name,
		ImporterEmail: pkg.Importer.Email,
		Importer:      pkg.Importer.Name,
	}

	return impl.genTemplate(impl.cfg.Template.AppendSigInfo, data)
}

func (impl *pullRequestImpl) genNewRepoData(pkg *domain.SoftwarePkg) (string, error) {
	data := struct {
		PkgName     string
		PkgDesc     string
		BranchName  string
		ProtectType string
		PublicType  string
	}{
		PkgName:     pkg.Name,
		PkgDesc:     pkg.Application.PackageDesc,
		BranchName:  impl.cfg.PR.NewRepoBranch.Name,
		ProtectType: impl.cfg.PR.NewRepoBranch.ProtectType,
		PublicType:  impl.cfg.PR.NewRepoBranch.PublicType,
	}

	return impl.genTemplate(impl.cfg.Template.NewRepoFile, data)
}

func (impl *pullRequestImpl) genTemplate(fileName string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (impl *pullRequestImpl) createPR(pkg *domain.SoftwarePkg) (pr sdk.PullRequest, err error) {
	prName := pkg.Name + impl.cfg.PR.PRName
	head := fmt.Sprintf("%s:%s", impl.robotLogin, impl.branchName(pkg.Name))
	return impl.cli.CreatePullRequest(
		impl.cfg.PR.Org, impl.cfg.PR.Repo, prName,
		pkg.Application.ReasonToImportPkg, head, "master", true,
	)
}
