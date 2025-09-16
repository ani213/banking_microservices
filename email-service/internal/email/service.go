package email

import "github.com/ani213/email-service/internal/config"

type Service struct {
	config *config.Config
}

func NewService(config *config.Config) *Service {
	return &Service{config: config}
}
