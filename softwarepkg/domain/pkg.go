package domain

import "time"

const (
	PkgStatusPRMerged    = "pr_merged"
	PkgStatusPRCreated   = "pr_created"
	PkgStatusPRException = "pr_exception"
	PkgStatusRepoCreated = "repo_created"
	PkgStatusInitialized = "initialized"
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

	Status            string
	PullRequest       PullRequest
	Importer          Importer
	Application       SoftwarePkgApplication
	CIPRNum           int
	PRExceptionExpiry int64
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

func (r *SoftwarePkg) SetPkgPRException() {
	r.Status = PkgStatusPRException
	r.PRExceptionExpiry = time.Now().Unix() + 600
}

func (r *SoftwarePkg) IsPRExceptionExpiried() bool {
	return time.Now().Unix() >= r.PRExceptionExpiry
}

func (r *SoftwarePkg) IsPkgStatusMerged() bool {
	return r.Status == PkgStatusPRMerged
}
