package domain

import "encoding/json"

const PlatformGitee = "gitee"

type PRCIFinishedEvent struct {
	PkgId        string `json:"pkg_id"`
	RelevantPR   string `json:"relevant_pr"`
	RepoLink     string `json:"repo_link"`
	FailedReason string `json:"failed_reason"`
}

func (e *PRCIFinishedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

type RepoCreatedEvent struct {
	PkgId        string `json:"pkg_id"`
	Platform     string `json:"platform"`
	RepoLink     string `json:"repo_link"`
	FailedReason string `json:"failed_reason"`
}

func (e *RepoCreatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

type CodePushedEvent = RepoCreatedEvent
