package main

import (
	"fmt"

	"os"
	"strings"

	wire "github.com/scanyourkube/cronjob/di"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	namespacesToIgnore := getEnvAsSlice("NAMESPACES_TO_IGNORE", []string{}, ",")
	vulnerabilitiesToIgnore := getEnvAsSlice("VULNERABILITIES_TO_IGNORE", []string{}, ",")
	kubeclarityApiUrl := os.Getenv("KUBECLARITY_API_URL")
	keelWebhookUrl := os.Getenv("KEEL_WEBHOOK_URL")
	senderEmailAddress := os.Getenv("SENDER_EMAIL")
	smtpServerAddress := os.Getenv("SMTP_SERVERADDRESS")
	smtpServerPort := os.Getenv("SMTP_SERVERPORT")

	fmt.Println(kubeclarityApiUrl)
	fmt.Println(keelWebhookUrl)

	kubernetesService := wire.InizializeKubernetesService(namespacesToIgnore)
	updateService := wire.InitializeKeelUpdateService(keelWebhookUrl)
	scanService := wire.InitializeScanService(kubeclarityApiUrl, kubernetesService, vulnerabilitiesToIgnore)
	kubeClarityResourceService := wire.InitializeResourceService(kubeclarityApiUrl, vulnerabilitiesToIgnore)
	notificationService := wire.InitializeNotificationService(senderEmailAddress, smtpServerAddress, smtpServerPort)
	automaticUpdateService := wire.InitializeAutomaticUpdateService(updateService, scanService, kubernetesService, notificationService, kubeClarityResourceService)

	err := automaticUpdateService.StartAutomaticUpdate()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Successfully performed automatic update. Check previous logs for details.")
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := os.Getenv(name)

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
