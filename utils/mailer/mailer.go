package mailer

import (
	"os"

	"github.com/mailjet/mailjet-apiv3-go/v3"
)

type MailService interface {
	SendResetPassword(to string, token string) error
}

type MailClient struct {
	mailClient *mailjet.Client
}

func NewMailService() MailService {
	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MJ_API_KEY"), os.Getenv("MJ_SECRET_KEY"))
	return &MailClient{
		mailClient: mailjetClient,
	}
}
