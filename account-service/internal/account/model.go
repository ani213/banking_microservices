package account

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID            int64           `db:"id" json:"id"`                                             // primary key
	UserID        int64           `db:"user_id" json:"userId" validate:"required"`                // FK -> users.id
	AccountNumber string          `db:"account_number" json:"accountNumber" validate:"required"`  // unique account number
	Balance       decimal.Decimal `db:"balance" json:"balance"`                                   // money (NUMERIC in DB)
	AccountTypeID int             `db:"account_type_id" json:"accountTypeId" validate:"required"` // FK -> account_types.id
	StatusID      int             `db:"status_id" json:"statusId" validate:"required"`            // FK -> account_status.id
	CreatedAt     time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time       `db:"updated_at" json:"updatedAt"`
}

type AccountRequest struct {
	UserId        int64  `json:"userId" validate:"required"`
	AccountNumber string `json:"accountNumber" validate:"required"`
	AccountTypeId int    `json:"accountTypeId" validate:"required"`
	Balance       int64  `json:"balance" validate:"required,min=0"`
	StatusID      int
}

type EmailRequestBody struct {
	To      string `json:"to"`
	Body    string `json:"body"`
	Subject string `json:"subject"`
}

type ContextValue struct {
	Email    string
	User_id  int64
	FullName string
}

type ContextKey string

var UserContextKey ContextKey = "user"

type ResponseAccount struct {
	UserID        int64           `db:"user_id" json:"userId" validate:"required"`                // FK -> users.id
	AccountNumber string          `db:"account_number" json:"accountNumber" validate:"required"`  // unique account number
	Balance       decimal.Decimal `db:"balance" json:"balance"`                                   // money (NUMERIC in DB)
	AccountTypeID int             `db:"account_type_id" json:"accountTypeId" validate:"required"` // FK -> account_types.id
	StatusID      int             `db:"status_id" json:"statusId" validate:"required"`            // FK -> account_status.id
	CreatedAt     time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time       `db:"updated_at" json:"updatedAt"`
	FullName      string          `db:"full_name" json:"fullName"`
	Email         string          `db:"email" json:"email"`
	AccountType   string          `db:"account_type" json:"account_type"`
	Status        string          `db:"status" json:"status"`
}

type ResponseAllUserWithAccount struct {
	UserId   string        `json:"userId"`
	Accounts []UserAccount `json:"accounts"`
}

type UserAccount struct {
	UserID        int64           `db:"user_id" json:"userId" `               // FK -> users.id
	AccountNumber string          `db:"account_number" json:"accountNumber" ` // unique account number
	Balance       decimal.Decimal `db:"balance" json:"balance"`               // money (NUMERIC in DB)
	// CreatedAt     time.Time       `db:"created_at" json:"createdAt"`
	// UpdatedAt     time.Time       `db:"updated_at" json:"updatedAt"`
	Email       string `db:"email" json:"email"`
	AccountType string `db:"account_type" json:"account_type"`
	Status      string `db:"status" json:"status"`
}
