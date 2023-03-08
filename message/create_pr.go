package message

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

const repoHandleScript = "./repo.sh"

type CreatePRParam domain.SoftwarePkgAppliedEvent

func (c CreatePRParam) modifyFiles(cfg *config) error {
	if err := c.appendToSigInfo(cfg); err != nil {
		return err
	}

	return c.newCreateRepoYaml(cfg)
}

func (c CreatePRParam) appendToSigInfo(cfg *config) error {
	appendContent, err := c.genAppendSigInfoData()
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf(cfg.PR.ModifyFiles.SigInfo, c.ImportingPkgSig)

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

func (c CreatePRParam) newCreateRepoYaml(cfg *config) error {
	subDirName := strings.ToLower(c.PkgName[:1])
	fileName := fmt.Sprintf(cfg.PR.ModifyFiles.NewRepo,
		c.ImportingPkgSig, subDirName, c.PkgName,
	)

	content, err := c.genNewRepoData(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, []byte(content), 0644)
}

type CmdType string

var (
	cmdInit      = CmdType("init")
	cmdNewBranch = CmdType("new")
	cmdCommit    = CmdType("commit")
)

func (c CreatePRParam) initRepo(cfg *config) error {
	if s, err := os.Stat(cfg.PR.Repo); err == nil && s.IsDir() {
		return nil
	}

	return c.execScript(cfg, cmdInit)
}

func (c CreatePRParam) newBranch(cfg *config) error {
	return c.execScript(cfg, cmdNewBranch)
}

func (c CreatePRParam) commit(cfg *config) error {
	return c.execScript(cfg, cmdCommit)
}

func (c CreatePRParam) execScript(cfg *config, cmdType CmdType) error {
	cmd := exec.Command(repoHandleScript, string(cmdType),
		cfg.Robot.Username, cfg.Robot.Password,
		cfg.Robot.Email, branchName(cfg.PR.BranchName, c.PkgName))

	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}

	return nil
}

func (c CreatePRParam) genAppendSigInfoData() (string, error) {
	data := struct {
		PkgName       string
		ImporterEmail string
		Importer      string
	}{
		PkgName:       c.PkgName,
		ImporterEmail: c.ImporterEmail,
		Importer:      c.Importer,
	}

	return genTemplate("./template/append_sig_info.tpl", data)
}

func (c CreatePRParam) genNewRepoData(cfg *config) (string, error) {
	data := struct {
		PkgName       string
		PkgDesc       string
		SourceCodeUrl string
		BranchName    string
		ProtectType   string
		PublicType    string
	}{
		PkgName:       c.PkgName,
		PkgDesc:       c.PkgDesc,
		SourceCodeUrl: c.SourceCodeURL,
		BranchName:    cfg.PR.NewRepoBranch.Name,
		ProtectType:   cfg.PR.NewRepoBranch.ProtectType,
		PublicType:    cfg.PR.NewRepoBranch.PublicType,
	}

	return genTemplate("./template/new_repo_file.tpl", data)
}

func genTemplate(fileName string, data interface{}) (string, error) {
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
