package types

type PlainTextEmail struct {
	To      string `json:"to" validate:"required"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body"`
}

type TemplateEmail struct {
	To      string         `json:"to" validate:"required"`
	Subject string         `json:"subject" validate:"required"`
	Data    map[string]any `json:"data"`
}

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
