package service_notification

import (
	"crypto/tls"
	"strconv"

	"gopkg.in/gomail.v2"

	log "github.com/sirupsen/logrus"
)

type EmailNotification struct {
	ToEmail string
	Subject string
	Body    string
}
type IEmailNotificationService interface {
	SendEmail(email EmailNotification) error
}

type EmailNotificationService struct {
	senderEmailAddress string
	host               string
	port               string
}

func NewEmailNotificationService(senderEmailAddress string, host string, port string) EmailNotificationService {
	return EmailNotificationService{
		senderEmailAddress: senderEmailAddress,
		host:               host,
		port:               port,
	}
}

func (service EmailNotificationService) SendEmail(email EmailNotification) error {
	port, err := strconv.Atoi(service.port)

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", service.senderEmailAddress)
	m.SetHeader("To", email.ToEmail)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", email.Body)

	d := gomail.NewDialer(service.host, port, service.senderEmailAddress, "")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	log.Debugf("Sending an email notification with subject %s to %s", email.Subject, email.ToEmail)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		log.Errorf("Error while sending mail to %s, error: %v", email.ToEmail, err)
		return err
	}

	return nil
}
