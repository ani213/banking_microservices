package auth

import (
	"errors"

	"github.com/ani213/banking-auth/pkg/jwtutil"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(email, password string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hash),
	}
	return s.repo.CreateUser(user)
}

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}
	return jwtutil.GenerateToken(user.ID)
}
