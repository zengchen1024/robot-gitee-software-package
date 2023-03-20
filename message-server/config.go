package messageserver

type Config struct {
	GroupName string `json:"group_name"    required:"true"`
	Topics    Topics `json:"topics"        required:"true"`
}

type Topics struct {
	NewPkg           string `json:"new_pkg"            required:"true"`
	CIPassed         string `json:"ci_passed"          required:"true"`
	ApprovedPkg      string `json:"approved_pkg"       required:"true"`
	RejectedPkg      string `json:"rejected_pkg"       required:"true"`
	AbandonedPkg     string `json:"abandoned_pkg"      required:"true"`
	CreatedRepo      string `json:"created_repo"       required:"true"`
	ClosedPR         string `json:"closed_pr"          required:"true"`
	MergedPR         string `json:"merged_pr"          required:"true"`
	AlreadyClosedPkg string `json:"already_closed_pkg" required:"true"`
}
