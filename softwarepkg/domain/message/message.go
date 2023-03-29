package message

type EventMessage interface {
	Message() ([]byte, error)
}

type SoftwarePkgMessage interface {
	NotifyCIResult(EventMessage) error
	NotifyRepoCreatedResult(EventMessage) error
	NotifyCodePushedResult(EventMessage) error
}
