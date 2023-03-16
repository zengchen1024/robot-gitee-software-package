package domain

import "encoding/json"

const platformGitee = "gitee"

type prCIFinishedEvent struct {
	PkgId        string `json:"pkg_id"`
	PkgName      string `json:"pkg_name"`
	RelevantPR   string `json:"relevant_pr"`
	FailedReason string `json:"failed_reason"`
	Success      bool   `json:"success"`
}

func (e *prCIFinishedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewPRCIFinishedEvent(pr *PullRequest, failedReason string) prCIFinishedEvent {
	return prCIFinishedEvent{
		PkgId:        pr.Pkg.Id,
		PkgName:      pr.Pkg.Name,
		RelevantPR:   pr.Link,
		FailedReason: failedReason,
		Success:      failedReason == "",
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

type prClosedEvent struct {
	PkgId      string `json:"pkg_id"`
	Reason     string `json:"reason"`
	RejectedBy string `json:"rejected_by"`
}

func (e *prClosedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewPRClosedEvent(pr *PullRequest, reason, rejectBy string) prClosedEvent {
	return prClosedEvent{
		PkgId:      pr.Pkg.Id,
		Reason:     reason,
		RejectedBy: rejectBy,
	}
}

type prMergedEvent struct {
	PkgId      string   `json:"pkg_id"`
	ApprovedBy []string `json:"approved_by"`
}

func (e *prMergedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewPRMergedEvent(pr *PullRequest, ApprovedBy []string) prMergedEvent {
	return prMergedEvent{
		PkgId:      pr.Pkg.Id,
		ApprovedBy: ApprovedBy,
	}
}
