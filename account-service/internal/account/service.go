package account

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ani213/account-service/internal/config"
	"github.com/ani213/account-service/internal/util"
	"github.com/shopspring/decimal"
)

type Service struct {
	repo   *Repository
	config *config.Config
}

func NewService(repo *Repository, config *config.Config) *Service {
	return &Service{repo: repo, config: config}
}

func (s *Service) CreateAccount(ctx context.Context, acc *Account) error {
	return s.repo.Create(ctx, acc)
}

func (s *Service) GetAccount(ctx context.Context, id int64) (*Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Deposit(ctx context.Context, id int64, amount decimal.Decimal) error {
	return s.repo.Deposit(ctx, id, amount)
}

func (s *Service) Withdraw(ctx context.Context, id int64, amount decimal.Decimal) error {
	return s.repo.Withdraw(ctx, id, amount)
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
	client := &http.Client{}
	resp, err := client.Do(emailReq)

	if err != nil {
		log.Println(err.Error(), "Error main things")
		return
	}
	log.Println(resp.StatusCode, "response")

	defer resp.Body.Close()
	log.Println("Email sent to:-  " + to)
}

func (s *Service) GetAccountsByUserID(userId string) ([]ResponseAccount, error) {
	accounts, err := s.repo.AccountsByUserID(userId)
	if err != nil {
		return []ResponseAccount{}, err
	}
	return accounts, nil
}

func (s *Service) GetAllUserWithAccounts() ([]UserAccount, error) {

	user, err := s.repo.GetAllUserWithAccounts()
	if err != nil {
		return []UserAccount{}, err
	}
	return user, nil

}
