package account

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID            int64           `db:"id" json:"id"`                         // primary key
	UserID        int64           `db:"user_id" json:"userId"`                // FK -> users.id
	AccountNumber string          `db:"account_number" json:"accountNumber"`  // unique account number
	Balance       decimal.Decimal `db:"balance" json:"balance"`               // money (NUMERIC in DB)
	AccountTypeID int             `db:"account_type_id" json:"accountTypeId"` // FK -> account_types.id
	StatusID      int             `db:"status_id" json:"statusId"`            // FK -> account_status.id
	CreatedAt     time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time       `db:"updated_at" json:"updatedAt"`
}
