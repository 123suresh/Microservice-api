package mailer

import (
	"fmt"
	"os"

	"github.com/mailjet/mailjet-apiv3-go/v3"
	"github.com/sirupsen/logrus"
)

func (mailClient *MailClient) SendResetPassword(to string, token string) error {
	logrus.Info("to => ", to)
	forgotUrl := fmt.Sprintf("%s/user/reset-password?token=%s", os.Getenv("HOST_URL"), token)
	messageInfo := []mailjet.InfoMessagesV31{
		{
			Priority: 1,
			From: &mailjet.RecipientV31{
				Email: "sureshthapa4009@gmail.com",
				Name:  "Suresh Thapa",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
					Name:  "User",
				},
			},
			Subject:  "Forget Password",
			TextPart: "Hi",
			HTMLPart: "<h3>Your forgot password link </h3>" + forgotUrl,
		},
	}
	messages := mailjet.MessagesV31{Info: messageInfo}
	res, err := mailClient.mailClient.SendMailV31(&messages)
	if err != nil {
		logrus.Error("Error to send email ", err)
		return err
	}
	logrus.Info("reset password mail successfully sent.", res)
	return nil
}
