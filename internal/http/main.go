package http

import (
	"log"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/mauriciofsnts/hermes/internal/config"
)

func Listen() error {
	smtpHost := config.Hermes.SmtpHost
	smtpPort := config.Hermes.SmtpPort
	smtpUsername := config.Hermes.SmtpUsername
	smtpPassword := config.Hermes.SmtpPassword

	// Endereço de e-mail padrão do remetente
	defaultFrom := config.Hermes.DefaultFrom

	http.HandleFunc("/send-email", func(w http.ResponseWriter, r *http.Request) {
		// Obtenha os dados do formulário POST
		to := r.FormValue("to")
		subject := r.FormValue("subject")
		body := r.FormValue("body")

		// Crie a mensagem de e-mail
		msg := []byte(
			"From: " + defaultFrom + "\r\n" +
				"To: " + to + "\r\n" +
				"Subject: " + subject + "\r\n" +
				"\r\n" +
				body + "\r\n",
		)

		// Autentique-se com o serviço SMTP
		auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

		// Envie a mensagem de e-mail
		err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, defaultFrom, []string{to}, msg)
		if err != nil {
			log.Fatal(err)
		}

		// Responda ao cliente com um status 200 OK
		w.WriteHeader(http.StatusOK)
	})

	// Inicie o servidor HTTP na porta 8080
	log.Fatal(http.ListenAndServe(":8080", nil))

	return nil
}
