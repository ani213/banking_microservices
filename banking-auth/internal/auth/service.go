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

func (s *Service) Register(userReq *RegisterRequest) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	user := &User{
		ID:           uuid.New().String(),
		Email:        userReq.Email,
		PasswordHash: string(hash),
		FullName:     userReq.FullName,
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

func (s *Service) GetUsers() ([]ResponsGetUser, error) {

	users, err := s.repo.FindUsers()

	if err != nil || users == nil {
		return nil, errors.New("invalid credentials")
	}
	return users, nil
}

func (s *Service) ValidateToken(token string) (string, error) {
	userID, err := jwtutil.ValidateToken(token)
	if err != nil {
		return "", err
	}
	return userID, nil
}
