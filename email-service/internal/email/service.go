package email

import (
	"context"
	"encoding/json"
	"log"
	"net/smtp"
	"time"

	"github.com/ani213/email-service/internal/config"
	"golang.org/x/time/rate"
)

type Service struct {
	config *config.Config
}

func NewService(config *config.Config) *Service {
	return &Service{config: config}
}

func (s *Service) SendEmail(requestBody *EmailRequest) (string, error) {
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
	if s.config.SMTPHost == "localhost" {
		auth = nil
	}
	err := smtp.SendMail(s.config.SMTPHost+":"+s.config.SMTPPort, auth, s.config.EmailId, to, msg)
	if err != nil {
		return "", err
	}
	return requestBody.To, nil
}

func (s *Service) StartRabbitConsumer() {
	log.Println("StartRabbitConsumer func called")
	ch := s.config.QueueChannel
	q, err := ch.QueueDeclare(
		"email-queue", //name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments

	)
	if err != nil {
		log.Println(err.Error(), "QueueDeclare error")
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println(err.Error(), "queue consume error")
	}

	// ✅ Rate limiter: 1 email every 2 seconds, burst up to 5
	limiter := rate.NewLimiter(rate.Every(12*time.Second), 5)

	log.Println("Consumer started. Waiting for messages...")

	for d := range msgs {
		var e EmailRequest
		if err := json.Unmarshal(d.Body, &e); err != nil {
			log.Println("Failed to parse message:", err)
			continue
		}
		// Wait for token before sending
		if err := limiter.Wait(context.Background()); err != nil {
			log.Println("Limiter error:", err)
			continue
		}
		// ✅ Send email via your service
		email, err := s.SendEmail(&e)
		if err != nil {
			log.Println("Failed to send email:", err)
		} else {
			log.Printf("Email sent to %s %s", email, time.Now())
		}
	}
}

func (s *Service) RetryEmailService() {
	log.Println("RetryEmailService func called")
	ch := s.config.QueueChannel
	q, err := ch.QueueDeclare(
		"email-queue", //name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments

	)
	if err != nil {
		log.Println(err.Error(), "QueueDeclare error")
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println(err.Error(), "queue consume error")
	}

	// ✅ Rate limiter: 1 email every 2 seconds, burst up to 5
	limiter := rate.NewLimiter(rate.Every(12*time.Second), 5)

	log.Println("Consumer started. Waiting for messages...")

	for d := range msgs {
		var e EmailRequest
		if err := json.Unmarshal(d.Body, &e); err != nil {
			log.Println("Failed to parse message:", err)
			continue
		}
		// Wait for token before sending
		if err := limiter.Wait(context.Background()); err != nil {
			log.Println("Limiter error:", err)
			continue
		}
		// ✅ Send email via your service
		email, err := s.SendEmail(&e)
		if err != nil {
			log.Println("Failed to send email:", err)
		} else {
			log.Printf("Email sent to %s %s", email, time.Now())
		}
	}
}
