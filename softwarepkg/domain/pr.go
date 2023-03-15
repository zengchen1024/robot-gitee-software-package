package domain

type SoftwarePkgSourceCode struct {
	Address string
	License string
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
	Num    int
	Link   string
	merged bool
	Pkg    SoftwarePkgBasic
}

func (r *PullRequest) SetMerged() {
	r.merged = true
}

func (r *PullRequest) IsMerged() bool {
	return r.merged
}

// SoftwarePkgRepo
type SoftwarePkgRepo struct {
	Pkg     SoftwarePkgBasic
	RepoURL string
}
