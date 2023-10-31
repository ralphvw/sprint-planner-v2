package helpers

import (
	"os"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendMail(content string, subject string, recepient string, recepientName string) {
	from := mail.NewEmail("Sprint Team", os.Getenv("SENDRID_MAIL"))
	to := mail.NewEmail(recepientName, recepient)
	message := mail.NewSingleEmail(from, subject, to, "", content)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		LogAction("EMAIL ERROR: " + err.Error())
	} else {
		LogAction("EMAIL SENT SUCCESSFULLY " + strconv.Itoa(response.StatusCode))
	}
}
