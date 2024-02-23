package mail

import (
	"fmt"
	"net/smtp"
	"strings"
	"task-one/helpers"
)

type Mailer interface {
	SendMail(to []string, cc []string, subject, message string) error
}

type SMTPMailer struct {
}

func (g *SMTPMailer) SendMail(to []string, cc []string, subject, message string) error {
	env := helpers.GetConfig()
	body := "From: " + env.Mail.SenderName + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", env.Mail.AuthEmail, env.Mail.AuthPassword, env.Mail.SmtpHost)
	smtpAddr := fmt.Sprintf("%s:%d", env.Mail.SmtpHost, env.Mail.SmtpPort)

	err := smtp.SendMail(smtpAddr, auth, env.Mail.AuthEmail, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil

}
