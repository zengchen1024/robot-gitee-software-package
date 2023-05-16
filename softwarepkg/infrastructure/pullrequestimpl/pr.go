package pullrequestimpl

import (
	"path/filepath"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

func NewPullRequestImpl(cfg *Config) (*pullRequestImpl, error) {
	localRepoDir, err := cloneRepo(cfg)
	if err != nil {
		return nil, err
	}

	cli := client.NewClient(func() []byte {
		return []byte(cfg.Robot.Token)
	})

	robot := client.NewClient(func() []byte {
		return []byte(cfg.CommunityRobot.Token)
	})

	tmpl, err := newtemplateImpl(&cfg.Template)
	if err != nil {
		return nil, err
	}

	return &pullRequestImpl{
		cli:          cli,
		cfg:          *cfg,
		template:     tmpl,
		cliToMergePR: robot,
		localRepoDir: localRepoDir,
	}, nil
}

func cloneRepo(cfg *Config) (string, error) {
	user := &cfg.Robot

	params := []string{
		cfg.ShellScript.CloneScript,
		cfg.ShellScript.WorkDir,
		user.Username,
		user.Email,
		user.Repo,
		user.cloneURL(),
		cfg.CommunityRobot.RepoLink,
	}

	if _, err, _ := utils.RunCmd(params...); err != nil {
		return "", err
	}

	return filepath.Join(cfg.ShellScript.WorkDir, user.Repo), nil
}

type iClient interface {
	CreatePullRequest(org, repo, title, body, head, base string, canModify bool) (sdk.PullRequest, error)
	GetGiteePullRequest(org, repo string, number int32) (sdk.PullRequest, error)
	ClosePR(org, repo string, number int32) error
	CreatePRComment(org, repo string, number int32, comment string) error
}

type clientToMergePR interface {
	MergePR(owner, repo string, number int32, opt sdk.PullRequestMergePutParam) error
}

type pullRequestImpl struct {
	cli          iClient
	cfg          Config
	template     templateImpl
	cliToMergePR clientToMergePR
	localRepoDir string
}

func (impl *pullRequestImpl) Create(pkg *domain.SoftwarePkg) (domain.PullRequest, error) {
	if err := impl.createBranch(pkg); err != nil {
		return domain.PullRequest{}, err
	}

	return impl.createPR(pkg)
}

func (impl *pullRequestImpl) Merge(prNum int) error {
	org := impl.cfg.CommunityRobot.Org
	repo := impl.cfg.CommunityRobot.Repo

	v, err := impl.cli.GetGiteePullRequest(org, repo, int32(prNum))
	if err != nil {
		return err
	}

	if v.State != sdk.StatusOpen {
		return nil
	}

	return impl.cliToMergePR.MergePR(
		org, repo, int32(prNum), sdk.PullRequestMergePutParam{},
	)
}

func (impl *pullRequestImpl) Close(prNum int) error {
	org := impl.cfg.CommunityRobot.Org
	repo := impl.cfg.CommunityRobot.Repo

	prDetail, err := impl.cli.GetGiteePullRequest(org, repo, int32(prNum))
	if err != nil {
		return err
	}

	if prDetail.State != sdk.StatusOpen {
		return nil
	}

	return impl.cli.ClosePR(org, repo, int32(prNum))
}

func (impl *pullRequestImpl) Comment(prNum int, content string) error {
	return impl.cli.CreatePRComment(
		impl.cfg.CommunityRobot.Org, impl.cfg.CommunityRobot.Repo,
		int32(prNum), content,
	)
}
