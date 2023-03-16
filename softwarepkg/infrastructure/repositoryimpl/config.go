package repositoryimpl

type Config struct {
	Table Table `json:"table" required:"true"`
}

type Table struct {
	SoftwarePkgPR string `json:"software_pkg_pr"    required:"true"`
}
