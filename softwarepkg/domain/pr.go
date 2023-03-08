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
	Application SoftwarePkgApplication
}

// PullRequest
type PullRequest struct {
	Num  int
	Link string
	Pkg  SoftwarePkgBasic
}

// SoftwarePkgRepo
type SoftwarePkgRepo struct {
	Pkg     SoftwarePkgBasic
	RepoURL string
}
