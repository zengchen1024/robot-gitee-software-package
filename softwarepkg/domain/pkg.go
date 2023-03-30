package domain

const (
	PkgStatusInitialized = "initialized"
	PkgStatusPRMerged    = "pr_merged"
	PkgStatusRepoCreated = "repo_created"
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
}

func (r *SoftwarePkg) SetPkgStatusInitialized() {
	r.Status = PkgStatusInitialized
}

func (r *SoftwarePkg) SetPkgStatusMerged() {
	r.Status = PkgStatusPRMerged
}

func (r *SoftwarePkg) SetPkgStatusRepoCreated() {
	r.Status = PkgStatusRepoCreated
}

func (r *SoftwarePkg) IsPkgStatusMerged() bool {
	return r.Status == PkgStatusPRMerged
}
