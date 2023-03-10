package email

type Email interface {
	Send(string) error
}
