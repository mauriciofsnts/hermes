package notification

import (
	"errors"
	"fmt"

	disgo "github.com/disgoorg/disgo/discord"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/discord"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/worker"
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type EmailController struct {
	Provider template.TemplateProvider
	Queue    worker.Queue[types.Mail]
}

func NewEmailController(template template.TemplateProvider, queue worker.Queue[types.Mail]) *EmailController {
	return &EmailController{
		Provider: template,
		Queue:    queue,
	}
}

var _ api.Controller = &EmailController{}

func (e *EmailController) Route(r api.Router) {
	r.Post("/notification", e.Notify)
}

func (e *EmailController) ValidateEmailNotification(templateId string, data map[string]any, subject string) (*types.Mail, error) {
	found := e.Provider.Exists(templateId)

	if !found {
		return nil, errors.New("template not found")
	}

	// Carregar template para validação
	templateContent, err := e.Provider.(*template.TemplateService).Get(templateId)
	if err != nil {
		return nil, errors.New("error loading template")
	}

	// Validar estrutura do template contra os dados fornecidos
	if err := template.ValidateTemplateStructure(templateContent, data); err != nil {
		return nil, err
	}

	// Renderizar template com validação de campos
	renderedTemplate, err := e.Provider.ParseHtmlTemplate(templateId, data)

	if err != nil {
		return nil, errors.New("error parsing template: " + err.Error())
	}

	if data["to"] == nil {
		return nil, errors.New("[to] field is required")
	}

	to, ok := data["to"].(string)

	if !ok {
		return nil, errors.New("recipient email is not a string")
	}

	notification := &types.Mail{
		To:      []string{to},
		Subject: subject,
		Sender:  config.Hermes.SMTP.Sender,
		Body:    renderedTemplate.String(),
	}

	return notification, nil
}

func (e *EmailController) ValidateDiscordNotification(apiKey string, data map[string]any, subject string) error {
	client, err := discord.Connect(apiKey)

	if err != nil {
		return errors.New(err.Error())
	}

	embed := disgo.NewEmbedBuilder().SetTitle(subject)

	for k, v := range data {
		strValue, ok := v.(string)
		if !ok {
			strValue = fmt.Sprintf("%v", v)
		}
		embed.AddField(k, strValue, false)
	}

	err = discord.SendWebhook(client, embed.Build())
	if err != nil {
		return errors.New("failed to send discord webhook: " + err.Error())
	}

	return nil
}
