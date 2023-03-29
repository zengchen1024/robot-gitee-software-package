package messageimpl

type Config struct {
	TopicsToNotify TopicsToNotify `json:"topics_to_notify"`
}

type TopicsToNotify struct {
	CIPassed    string `json:"ci_passed"          required:"true"`
	CreatedRepo string `json:"created_repo"       required:"true"`
	PushedCode  string `json:"pushed_code"        required:"true"`
}
