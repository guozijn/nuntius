package nuntius

type Provider interface {
	Send(message Message) error
}

type Message struct {
	To   []string
	From string
	Text string
}
