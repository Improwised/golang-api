package workers

import (
	"log"

	helpers "github.com/Improwised/golang-api/helpers/smtp"
)

type WelcomeMail struct {
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Roles     string `json:"roles"`
}

func (w WelcomeMail) Handle() error {
	log.Println("sending mail")

	smtp := helpers.NewSMTPHelper("sandbox.smtp.mailtrap.io", "2525", "7628fa366c0257c", "2afb7200812272")
	smtp.SetSubject("welocme")
	smtp.SetPlainBody([]byte("welcome to our org"))

	smtp.SetSender("chintansakhiya00001@gmail.com")
	smtp.SetReceivers([]string{"chintansakhiya00001@gmail.com"})

	if err := smtp.SendMail(); err != nil {
		return err
	}

	log.Printf("mail send to %v", w.Email)
	return nil
}
