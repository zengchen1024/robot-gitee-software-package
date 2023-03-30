package domain

const (
	StatusInitialized = "initialized"
	StatusPRMerged    = "pr_merged"
	StatusRepoCreated = "repo_created"
)

type SoftwarePkgSourceCode struct {
	SpecURL   string
	SrcRPMURL string
}

type SoftwarePkgApplication struct {
	SourceCode        SoftwarePkgSourceCode
	PackageDesc       string
	PackagePlatform   string
	ImportingPkgSig   string
	ReasonToImportPkg string
}

type SoftwarePkgBasic struct {
	Id   string
	Name string
}

type SoftwarePkg struct {
	SoftwarePkgBasic

	ImporterName  string
	ImporterEmail string
	Application   SoftwarePkgApplication
}

// PullRequest
type PullRequest struct {
	Num           int
	Link          string
	Status        string
	ImporterName  string
	ImporterEmail string
	Pkg           SoftwarePkgBasic
	SrcCode       SoftwarePkgSourceCode
}

func (r *PullRequest) SetStatusInitialized() {
	r.Status = StatusInitialized
}

func (r *PullRequest) SetStatusMerged() {
	r.Status = StatusPRMerged
}

func (r *PullRequest) SetStatusRepoCreated() {
	r.Status = StatusRepoCreated
}

func (r *PullRequest) IsStatusMerged() bool {
	return r.Status == StatusPRMerged
}

// SoftwarePkgRepo
type SoftwarePkgRepo struct {
	Pkg     SoftwarePkgBasic
	RepoURL string
}
