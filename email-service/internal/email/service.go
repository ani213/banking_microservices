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
	log.Panicln("Email Service called:- " + requestBody.To)
	to := []string{requestBody.To}
	msg := []byte("Subject: " + requestBody.Subject + "\n" + requestBody.Body)
	auth := smtp.PlainAuth("", s.config.EmailId, s.config.EmailPassword, s.config.SMTPHost)
	err := smtp.SendMail(s.config.SMTPHost+":"+s.config.SMTPPort, auth, s.config.EmailId, to, msg)
	if err != nil {
		return "", err
	}
	return "Email successfully sent" + requestBody.To, nil
}
