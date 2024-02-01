package workers

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type WelcomeMail struct {
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Roles     string `json:"roles"`
}

func (w WelcomeMail) Handle() error {
	log.Println("sending mail")

	m := gomail.NewMessage()
	m.SetHeader("From", "support@improwised.com")
	m.SetHeader("To", w.Email)
	m.SetHeader("Subject", "Welcome")
	m.SetBody("text/html", "welcome to our org")
	m.Attach("./../README.md")
	 
	d:=gomail.NewDialer("smtp.mailtrap.io", 2525, "7628fa366c0257", "2afb7200812272")
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	log.Printf("mail send to %v", w.Email)
	return nil
}

type LoginMail struct {
	Email string `json:"email"`
	Device string `json:"device"`
}

func (l LoginMail) Handle() error {
	fmt.Println(l)
	log.Println("login sending mail")

	m := gomail.NewMessage()
	m.SetHeader("From", "support@improwised.com")
	m.SetHeader("To", l.Email)
	m.SetHeader("Subject", "Login")
	m.SetBody("text/html", "login from "+l.Device)

	d:=gomail.NewDialer("smtp.mailtrap.io", 2525, "7628fa366c0257", "2afb7200812272")
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	log.Printf("mail send to %v", l.Email)
	return nil

}