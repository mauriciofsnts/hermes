package types

type MailType int

const (
	HTML MailType = 1
	TEXT MailType = 2
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
	Type    MailType
}
