package message

type EventMessage interface {
	Message() ([]byte, error)
}

type SoftwarePkgMessage interface {
	NotifyCIResult(EventMessage) error
	NotifyRepoCreatedResult(EventMessage) error
	NotifyPRClosed(EventMessage) error
	NotifyPRMerged(EventMessage) error
}
