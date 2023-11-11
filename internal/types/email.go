package types

type Email struct {
	To           string         `json:"to" validate:"required"`
	Subject      string         `json:"subject" validate:"required"`
	Body         string         `json:"body"`
	Content      map[string]any `json:"content"`
	TemplateName string         `json:"templateName"`
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
