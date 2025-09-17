package email

type EmailRequest struct {
	To      string `json:"to"`
	Body    string `json:"body"`
	Subject string `json:"subject"`
}
