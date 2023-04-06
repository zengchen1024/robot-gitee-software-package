package community

type Config struct {
	PkgOrg         string `json:"pkg_org"          required:"true"`
	CommunityOrg   string `json:"community_org"    required:"true"`
	CommunityRepo  string `json:"community_repo"   required:"true"`
	CISuccessLabel string `json:"ci_success_label" required:"true"`
	CIFailureLabel string `json:"ci_failure_label" required:"true"`
}

func (cfg *Config) isCommunity(org, repo string) bool {
	return cfg.CommunityOrg == org && cfg.CommunityRepo == repo
}
