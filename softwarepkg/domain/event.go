package domain

import "encoding/json"

type prCIFinishedEvent struct {
	PkgId   string `json:"pkg_id"`
	PkgName string `json:"pkg_name"`
	PRLink  string `json:"pr_link"`
}

func (e *prCIFinishedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewPRCIFinishedEvent(pr *PullRequest) prCIFinishedEvent {
	return prCIFinishedEvent{}
}
