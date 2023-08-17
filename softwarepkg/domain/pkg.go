package domain

const (
	PkgStatusInitialized = "initialized"
	PkgStatusPRCreated   = "pr_created"
	PkgStatusPRMerged    = "pr_merged"
	PkgStatusRepoCreated = "repo_created"
	PkgStatusException   = "exception" // more information in the email of maintainer
)

type SoftwarePkgSourceCode struct {
	SpecURL   string
	SrcRPMURL string
}

type SoftwarePkgApplication struct {
	Upstream          string
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

type PullRequest struct {
	Num  int
	Link string
}

type Importer struct {
	Name  string
	Email string
}

type SoftwarePkg struct {
	SoftwarePkgBasic

	Status      string
	PullRequest PullRequest
	Importer    Importer
	Application SoftwarePkgApplication
	CIPRNum     int
}

func (r *SoftwarePkg) SetPkgStatusInitialized() {
	r.Status = PkgStatusInitialized
}

func (r *SoftwarePkg) SetPkgStatusPRCreated() {
	r.Status = PkgStatusPRCreated
}

func (r *SoftwarePkg) SetPkgStatusPRMerged() {
	r.Status = PkgStatusPRMerged
}

func (r *SoftwarePkg) SetPkgStatusRepoCreated() {
	r.Status = PkgStatusRepoCreated
}

func (r *SoftwarePkg) SetPkgStatusException() {
	r.Status = PkgStatusException
}

func (r *SoftwarePkg) IsPkgStatusMerged() bool {
	return r.Status == PkgStatusPRMerged
}
