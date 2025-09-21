package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ani213/banking-auth/internal/config"
	"github.com/ani213/banking-auth/pkg/jwtutil"
	"github.com/ani213/banking-auth/util"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   *Repository
	config *config.Config
}

func NewService(repo *Repository, config *config.Config) *Service {
	return &Service{repo: repo, config: config}
}

func (s *Service) Register(userReq *RegisterRequest) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	user := &User{
		Email:        userReq.Email,
		PasswordHash: string(hash),
		FullName:     userReq.FullName,
	}
	return s.repo.CreateUser(user)
}

func (s *Service) Login(email, password string) (string, error) {

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("emailfind invalid credentials")
	}
	roles, err := s.repo.FindRoleByUserId((user.ID))
	if err != nil || user == nil {
		return "", err
	}
	// fmt.Println(roles, "roless")
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", errors.New("password invalid credentials")
	}
	return jwtutil.GenerateToken(user.ID, user.Email, user.FullName, roles)
}

func (s *Service) GetUsers() ([]ResponsGetUser, error) {

	users, err := s.repo.FindUsers()

	if err != nil || users == nil {
		return nil, errors.New("invalid credentials")
	}
	return users, nil
}

func (s *Service) ValidateToken(r *http.Request) (jwtutil.ContextValue, error) {
	user, err := jwtutil.GetContextValue(r)
	if err != nil {
		return jwtutil.ContextValue{}, err
	}
	return user, err

}

func (s *Service) SendEmail(to string, subject string, body string, r *http.Request) {
	token := util.GetToken(r)
	req := EmailRequestBody{
		To:      to,
		Body:    body,
		Subject: subject,
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	emailReq, err := http.NewRequest("POST", s.config.EmailServer+"/send-email", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println(err.Error())
	}
	emailReq.Header.Set("Content-Type", "application/json")
	emailReq.Header.Set("Authorization", "Bearer "+token)
	log.Println(emailReq.Header.Get("Authorization"))
	client := &http.Client{}
	resp, err := client.Do(emailReq)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(resp.StatusCode, "response")
	defer resp.Body.Close()
	log.Println("Email sent to:-  " + to)
}

func (s *Service) AddRoles(userId int64, roles []int64) error {
	return s.repo.AddRoles(userId, roles)
}
