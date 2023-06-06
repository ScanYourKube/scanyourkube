package service_notification

import (
	"fmt"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
)

func TestRunningMailServer(t *testing.T) {

	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})

	// To start server use Start() method
	if err := server.Start(); err != nil {
		fmt.Println(err)
	}

	hostAddress, portNumber := "127.0.0.1", server.PortNumber()

	emailNotificationService := NewEmailNotificationService("test@test.com", hostAddress, fmt.Sprintf("%d", portNumber))

	emailNoficiation := EmailNotification{
		ToEmail: "dummy@dummy.com",
		Subject: "Test",
		Body:    "Test",
	}
	err := emailNotificationService.SendEmail(emailNoficiation)

	if err != nil {
		t.Error(err)
	}
}

func TestInActiveMailServer(t *testing.T) {

	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})

	// To start server use Start() method
	if err := server.Start(); err != nil {
		fmt.Println(err)
	}

	hostAddress, portNumber := "127.0.0.1", server.PortNumber()

	if err := server.Stop(); err != nil {
		fmt.Println(err)
	}

	emailNotificationService := NewEmailNotificationService("test@test.com", hostAddress, fmt.Sprintf("%d", portNumber))

	emailNoficiation := EmailNotification{
		ToEmail: "dummy@dummy.com",
		Subject: "Test",
		Body:    "Test",
	}
	err := emailNotificationService.SendEmail(emailNoficiation)

	if err == nil {
		t.Error("Expected error")
	}
}

func TestIrregularPortNumber(t *testing.T) {

	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})

	// To start server use Start() method
	if err := server.Start(); err != nil {
		fmt.Println(err)
	}

	hostAddress, portNumber := "127.0.0.1", ""

	emailNotificationService := NewEmailNotificationService("test@test.com", hostAddress, portNumber)

	emailNoficiation := EmailNotification{
		ToEmail: "dummy@dummy.com",
		Subject: "Test",
		Body:    "Test",
	}
	err := emailNotificationService.SendEmail(emailNoficiation)

	if err == nil {
		t.Error("Expected error")
	}
}
