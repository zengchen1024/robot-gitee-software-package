package pullrequestimpl

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

type CmdType string

var (
	cmdInit      = CmdType("init")
	cmdNewBranch = CmdType("new")
	cmdCommit    = CmdType("commit")
)

func (impl *pullRequestImpl) initRepo() error {
	if s, err := os.Stat(impl.cfg.PR.Repo); err == nil && s.IsDir() {
		return nil
	}

	return impl.execScript(cmdInit)
}

func (impl *pullRequestImpl) newBranch() error {
	return impl.execScript(cmdNewBranch)
}

func (impl *pullRequestImpl) commit() error {
	return impl.execScript(cmdCommit)
}

func (impl *pullRequestImpl) execScript(cmdType CmdType) error {
	cmd := exec.Command(impl.cfg.ShellScript, string(cmdType),
		impl.cfg.Robot.Username, impl.cfg.Robot.Password,
		impl.cfg.Robot.Email, impl.branchName())

	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}

	return nil
}

func (impl *pullRequestImpl) modifyFiles() error {
	if err := impl.appendToSigInfo(); err != nil {
		return err
	}

	return impl.newCreateRepoYaml()
}

func (impl *pullRequestImpl) appendToSigInfo() error {
	appendContent, err := impl.genAppendSigInfoData()
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("community/sig/%s/sig-info.yaml",
		impl.pkg.Application.ImportingPkgSig)

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

func (impl *pullRequestImpl) newCreateRepoYaml() error {
	subDirName := strings.ToLower(impl.pkg.Name[:1])
	fileName := fmt.Sprintf("community/sig/%s/src-openeuler/%s/%s.yaml",
		impl.pkg.Application.ImportingPkgSig, subDirName, impl.pkg.Name)

	content, err := impl.genNewRepoData()
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, []byte(content), 0644)
}

func (impl *pullRequestImpl) branchName() string {
	return fmt.Sprintf("software_package_%s", impl.pkg.Name)
}

func (impl *pullRequestImpl) prName() string {
	return impl.pkg.Name + impl.cfg.PR.PRName
}

func (impl *pullRequestImpl) genAppendSigInfoData() (string, error) {
	data := struct {
		PkgName       string
		ImporterEmail string
		Importer      string
	}{
		PkgName:       impl.pkg.Name,
		ImporterEmail: impl.pkg.ImporterEmail,
		Importer:      impl.pkg.ImporterName,
	}

	return impl.genTemplate("./template/append_sig_info.tpl", data)
}

func (impl *pullRequestImpl) genNewRepoData() (string, error) {
	data := struct {
		PkgName       string
		PkgDesc       string
		SourceCodeUrl string
		BranchName    string
		ProtectType   string
		PublicType    string
	}{
		PkgName:       impl.pkg.Name,
		PkgDesc:       impl.pkg.Application.PackageDesc,
		SourceCodeUrl: impl.pkg.Application.SourceCode.Address,
		BranchName:    impl.cfg.PR.NewRepoBranch.Name,
		ProtectType:   impl.cfg.PR.NewRepoBranch.ProtectType,
		PublicType:    impl.cfg.PR.NewRepoBranch.PublicType,
	}

	return impl.genTemplate("./template/new_repo_file.tpl", data)
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

func (impl *pullRequestImpl) getRobotLogin() (string, error) {
	if impl.robotLogin == "" {
		v, err := impl.cli.GetBot()
		if err != nil {
			return "", err
		}

		impl.robotLogin = v.Login
	}

	return impl.robotLogin, nil
}

func (impl *pullRequestImpl) submit() (dpr domain.PullRequest, err error) {
	robotName, err := impl.getRobotLogin()
	if err != nil {
		return
	}

	head := fmt.Sprintf("%s:%s", robotName, impl.branchName())
	pr, err := impl.cli.CreatePullRequest(
		impl.cfg.PR.Org, impl.cfg.PR.Repo, impl.prName(),
		impl.pkg.Application.ReasonToImportPkg, head, "master", true,
	)
	if err != nil {
		return
	}

	dpr = domain.PullRequest{
		Num:  int(pr.Number),
		Link: pr.Url,
		Pkg:  impl.pkg.SoftwarePkgBasic,
	}

	return
}
