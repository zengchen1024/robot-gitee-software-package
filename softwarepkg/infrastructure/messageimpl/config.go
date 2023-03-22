package messageimpl

type Config struct {
	TopicsToNotify TopicsToNotify `json:"topics_to_notify"`
}

type TopicsToNotify struct {
	CIPassed    string `json:"ci_passed"          required:"true"`
	CreatedRepo string `json:"created_repo"       required:"true"`
	ClosedPR    string `json:"closed_pr"          required:"true"`
	MergedPR    string `json:"merged_pr"          required:"true"`
}
