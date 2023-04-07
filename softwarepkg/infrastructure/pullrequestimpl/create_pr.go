package pullrequestimpl

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

func (impl *pullRequestImpl) createBranch(pkg *domain.SoftwarePkg) error {
	sigInfoData, err := impl.genAppendSigInfoData(pkg)
	if err != nil {
		return err
	}

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
		fmt.Sprintf("sig/%s/sig-info.yaml", pkg.Application.ImportingPkgSig),
		sigInfoData,
		fmt.Sprintf(
			"sig/%s/src-openeuler/%s/%s.yaml",
			pkg.Application.ImportingPkgSig,
			strings.ToLower(pkg.Name[:1]),
			pkg.Name,
		),
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
	return impl.template.genSigInfo(&sigInfoTplData{
		PkgName:       pkg.Name,
		ImporterEmail: pkg.Importer.Email,
		Importer:      pkg.Importer.Name,
	})

}

func (impl *pullRequestImpl) genNewRepoData(pkg *domain.SoftwarePkg) (string, error) {
	return impl.template.genRepoYaml(&repoYamlTplData{
		PkgName:     pkg.Name,
		PkgDesc:     pkg.Application.PackageDesc,
		BranchName:  impl.cfg.PR.NewRepoBranch.Name,
		ProtectType: impl.cfg.PR.NewRepoBranch.ProtectType,
		PublicType:  impl.cfg.PR.NewRepoBranch.PublicType,
	})
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

func (impl *pullRequestImpl) createPR(pkg *domain.SoftwarePkg) (pr domain.PullRequest, err error) {
	body, err := impl.template.genPRBody(&prBodyTplData{
		PkgName: pkg.Name,
		PkgLink: impl.cfg.SoftwarePkg.Endpoint + pkg.Id,
	})
	if err != nil {
		return
	}

	v, err := impl.cli.CreatePullRequest(
		impl.cfg.PR.Org, impl.cfg.PR.Repo,
		fmt.Sprintf("add eco-package: %s", pkg.Name),
		body,
		fmt.Sprintf(
			"%s:%s", impl.cfg.Robot.Username, impl.branchName(pkg.Name),
		),
		"master", true,
	)
	if err == nil {
		pr.Num = int(v.Number)
		pr.Link = v.HtmlUrl
	}

	return
}
