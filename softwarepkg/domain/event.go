package domain

import "encoding/json"

const platformGitee = "gitee"

type prCIFinishedEvent struct {
	PkgId        string `json:"pkg_id"`
	RelevantPR   string `json:"relevant_pr"`
	RepoLink     string `json:"repo_link"`
	FailedReason string `json:"failed_reason"`
}

func (e *prCIFinishedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewPRCIFinishedEvent(
	pr *PullRequest, failedReason, repoLink string,
) prCIFinishedEvent {
	return prCIFinishedEvent{
		PkgId:        pr.Pkg.Id,
		RelevantPR:   pr.Link,
		RepoLink:     repoLink,
		FailedReason: failedReason,
	}
}

type repoCreatedEvent struct {
	PkgId    string `json:"pkg_id"`
	Platform string `json:"platform"`
	RepoLink string `json:"repo_link"`
}

func (e *repoCreatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewRepoCreatedEvent(pr *PullRequest, url string) repoCreatedEvent {
	return repoCreatedEvent{
		PkgId:    pr.Pkg.Id,
		Platform: platformGitee,
		RepoLink: url,
	}
}
