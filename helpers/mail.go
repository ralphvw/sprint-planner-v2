package helpers

import (
	"os"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(content string, subject string, recepient string, recepientName string) error {
	from := mail.NewEmail("Sprint Team", os.Getenv("SENDGRID_MAIL"))
	to := mail.NewEmail(recepientName, recepient)
	message := mail.NewSingleEmail(from, subject, to, "", content)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		LogAction("EMAIL ERROR: " + err.Error())
		return err
	} else {
		LogAction("EMAIL SENT SUCCESSFULLY " + strconv.Itoa(response.StatusCode) + " RESPONSE: " + response.Body)
	}
	return nil
}
