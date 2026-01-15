package account

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ani213/account-service/internal/config"
	"github.com/ani213/account-service/internal/util"
	amqp "github.com/rabbitmq/amqp091-go"
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

func (s *Service) Deposit(ctx context.Context, accountNo string, amount float64) (BalancEmail, error) {
	return s.repo.Deposit(ctx, accountNo, amount)
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

func (s *Service) GetAllUserWithAccounts() ([]ResponseAllUserWithAccount, error) {

	user, err := s.repo.GetAllUserWithAccounts()
	if err != nil {
		return []ResponseAllUserWithAccount{}, err
	}
	group := make(map[int64][]UserAccount)
	for _, item := range user {
		group[item.UserID] = append(group[item.UserID], UserAccount{AccountNumber: item.AccountNumber, Balance: item.Balance, AccountType: item.AccountType, UserID: item.UserID, Status: item.Status, Email: item.Email})
	}
	var result []ResponseAllUserWithAccount
	for key, value := range group {
		result = append(result, ResponseAllUserWithAccount{UserId: key, Email: value[0].Email, Accounts: value})
	}
	return result, nil

}

func (s *Service) GetEmaiByUserId(userId string) (string, error) {
	email, err := s.repo.GetEmailByUserId(userId)
	if err != nil {
		log.Println(err.Error(), "Error in service")
		return "", err
	}
	log.Println(email, "Email in service")
	return email, nil
}

func (s *Service) SendEmailInQueue(email EmailRequestBody) {
	var ch = s.config.QueueChannel
	q, err := ch.QueueDeclare(
		"email-queue", //name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments

	)
	if err != nil {
		log.Println("Error during declear queue", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body, err := json.Marshal(email)
	if err != nil {
		log.Println("Error durin json Marshal")
	}
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Println("Failed to publish a message", err.Error())
	}
}
