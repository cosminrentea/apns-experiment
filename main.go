package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"github.com/smancke/guble/protocol"
	"github.com/smancke/guble/server"
	"os"
)

const (
	defaultCertFileName = "development-certificate.p12"
)

func init() {
	server.AfterMessageDelivery = func(m *protocol.Message) {
		fmt.Print("message delivered")
	}
}

type APNSConfig struct {
	CertFileName string
	CertPassword string
	Topic        string
}

func main() {

	// server.Main()

	cfg := APNSConfig{
		CertFileName: defaultCertFileName,
		CertPassword: os.Getenv("APNS_CERT_PASSWORD"),
		Topic:        os.Getenv("APNS_TOPIC"),
	}
	deviceToken := os.Getenv("APNS_DEVICE_TOKEN")

	p := payload.NewPayload().
		AlertTitle("REWE Sonderrabatt").
		AlertBody("Sie haben ein Sonderrabatt von 50% für das neue iPhone 8 bekommen!").
		ContentAvailable()

	sendAPNSNotification(cfg, deviceToken, p)
}

func sendAPNSNotification(c APNSConfig, deviceToken string, p *payload.Payload) {
	cert, errCert := certificate.FromP12File(c.CertFileName, c.CertPassword)
	if errCert != nil {
		log.WithError(errCert).Error("APNS certificate error")
	}

	notification := &apns.Notification{}
	notification.Priority = apns.PriorityHigh
	notification.Topic = c.Topic
	notification.DeviceToken = deviceToken
	notification.Payload = p

	client := apns.NewClient(cert).Development()
	response, errPush := client.Push(notification)
	if errPush != nil {
		log.WithError(errPush).Error("APNS error when pushing notification")
		return
	}
	if response.Sent() {
		log.WithField("id", response.ApnsID).Debug("APNS notification successfully sent")
	} else {
		log.WithField("id", response.ApnsID).Error("APNS notification was not sent")
	}
}
