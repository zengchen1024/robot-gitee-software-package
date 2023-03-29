package pullrequestimpl

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

type CmdType string

var (
	cmdInit      = CmdType("init")
	cmdNewBranch = CmdType("new")
	cmdCommit    = CmdType("commit")
)

func (impl *pullRequestImpl) initRepo(pkg *domain.SoftwarePkg) error {
	if s, err := os.Stat(impl.cfg.PR.Repo); err == nil && s.IsDir() {
		return nil
	}

	return impl.execScript(cmdInit, pkg.Name)
}

func (impl *pullRequestImpl) newBranch(pkg *domain.SoftwarePkg) error {
	return impl.execScript(cmdNewBranch, pkg.Name)
}

func (impl *pullRequestImpl) commit(pkg *domain.SoftwarePkg) error {
	return impl.execScript(cmdCommit, pkg.Name)
}

func (impl *pullRequestImpl) execScript(cmdType CmdType, pkgName string) error {
	cmd := exec.Command(impl.cfg.ShellScript, string(cmdType),
		impl.cfg.Robot.Username, impl.cfg.Robot.Token,
		impl.cfg.Robot.Email, impl.branchName(pkgName))

	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}

	return nil
}

func (impl *pullRequestImpl) modifyFiles(pkg *domain.SoftwarePkg) error {
	if err := impl.appendToSigInfo(pkg); err != nil {
		return err
	}

	return impl.newCreateRepoYaml(pkg)
}

func (impl *pullRequestImpl) appendToSigInfo(pkg *domain.SoftwarePkg) error {
	appendContent, err := impl.genAppendSigInfoData(pkg)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf(
		"community/sig/%s/sig-info.yaml",
		pkg.Application.ImportingPkgSig)

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	if _, err = write.WriteString(appendContent); err != nil {
		return err
	}

	if err = write.Flush(); err != nil {
		return err
	}

	return nil
}

func (impl *pullRequestImpl) newCreateRepoYaml(pkg *domain.SoftwarePkg) error {
	subDirName := strings.ToLower(pkg.Name[:1])
	fileName := fmt.Sprintf(
		"community/sig/%s/src-openeuler/%s/%s.yaml",
		pkg.Application.ImportingPkgSig, subDirName, pkg.Name)

	content, err := impl.genNewRepoData(pkg)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
		return err
	}

	return os.WriteFile(fileName, []byte(content), 0644)
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
		ImporterEmail: pkg.ImporterEmail,
		Importer:      pkg.ImporterName,
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

func (impl *pullRequestImpl) submit(pkg *domain.SoftwarePkg) (pr sdk.PullRequest, err error) {
	prName := pkg.Name + impl.cfg.PR.PRName
	head := fmt.Sprintf("%s:%s", impl.robotLogin, impl.branchName(pkg.Name))
	return impl.cli.CreatePullRequest(
		impl.cfg.PR.Org, impl.cfg.PR.Repo, prName,
		pkg.Application.ReasonToImportPkg, head, "master", true,
	)
}

func (impl *pullRequestImpl) toPullRequest(
	pr *sdk.PullRequest, pkg *domain.SoftwarePkg,
) domain.PullRequest {
	return domain.PullRequest{
		Num:           int(pr.Number),
		Link:          pr.HtmlUrl,
		Pkg:           pkg.SoftwarePkgBasic,
		ImporterName:  pkg.ImporterName,
		ImporterEmail: pkg.ImporterEmail,
		SrcCode: domain.SoftwarePkgSourceCode{
			SpecURL:   pkg.Application.SourceCode.SpecURL,
			SrcRPMURL: pkg.Application.SourceCode.SrcRPMURL,
		},
	}
}
