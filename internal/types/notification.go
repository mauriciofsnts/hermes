package types

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

type NotificationRequest struct {
	TemplateID string      `json:"templateId"`
	Subject    string      `json:"subject"`
	Recipients []Recipient `json:"recipients"`
}

const (
	MAIL    RecipientType = "mail"
	DISCORD RecipientType = "discord"
)

type RecipientType string

type Recipient struct {
	Type RecipientType  `json:"type"`
	Data map[string]any `json:"data"`
}
