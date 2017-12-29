package notify

type Sender interface {
	Send([]byte) error
}
