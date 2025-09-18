package email

import (
	"log"
	"net/smtp"

	"github.com/ani213/email-service/internal/config"
)

type Service struct {
	config *config.Config
}

func NewService(config *config.Config) *Service {
	return &Service{config: config}
}

func (s *Service) SendEmail(requestBody *EmailRequest) (string, error) {
	log.Println("Email Service called:- "+requestBody.To, requestBody.Subject, requestBody.Body)
	to := []string{requestBody.To}
	// msg := []byte("Subject: " + requestBody.Subject + "\n" + requestBody.Body)
	msg := []byte(
		"Subject: " + requestBody.Subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
			"\r\n" + // <-- separates headers and body
			requestBody.Body,
	)
	auth := smtp.PlainAuth("", s.config.EmailId, s.config.EmailPassword, s.config.SMTPHost)
	err := smtp.SendMail(s.config.SMTPHost+":"+s.config.SMTPPort, auth, s.config.EmailId, to, msg)
	if err != nil {
		return "", err
	}
	return "Email successfully sent to:- " + requestBody.To + " body is:- " + requestBody.Body, nil
}
