package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"github.com/smancke/guble/protocol"
	"github.com/smancke/guble/server"
	"os"
)

const (
	defaultCertFileName = "development-certificate.p12"
	msgAPNSNotSent      = "APNS notification was not sent"
)

var (
	errAPNSNotSent = errors.New(msgAPNSNotSent)
)

func init() {
	server.AfterMessageDelivery = func(m *protocol.Message) {
		fmt.Print("message delivered")
	}
}

func main() {

	// server.Main()

	topic := os.Getenv("APNS_TOPIC")
	deviceToken := os.Getenv("APNS_DEVICE_TOKEN")
	cl := getAPNSClient(defaultCertFileName, os.Getenv("APNS_CERT_PASSWORD"), false)
	p := payload.NewPayload().
		AlertTitle("Sonderrabatt").
		AlertBody("Sie haben ein Sonderrabatt von 60% für das neue iPhone 8 erhalten!").
		ContentAvailable()
	sendAPNSNotification(cl, topic, deviceToken, p)
}

func getAPNSClient(certFileName string, certPassword string, production bool) *apns2.Client {
	cert, errCert := certificate.FromP12File(certFileName, certPassword)
	if errCert != nil {
		log.WithError(errCert).Error("APNS certificate error")
	}
	if production {
		return apns2.NewClient(cert).Production()
	}
	return apns2.NewClient(cert).Development()
}

func sendAPNSNotification(cl *apns2.Client, topic string, deviceToken string, p *payload.Payload) error {
	notification := &apns2.Notification{
		Priority:    apns2.PriorityHigh,
		Topic:       topic,
		DeviceToken: deviceToken,
		Payload:     p,
	}
	response, errPush := cl.Push(notification)
	if errPush != nil {
		log.WithError(errPush).Error("APNS error when trying to push notification")
		return errPush
	}
	if !response.Sent() {
		log.WithField("id", response.ApnsID).WithField("reason", response.Reason).Error(msgAPNSNotSent)
		return errAPNSNotSent
	}
	log.WithField("id", response.ApnsID).Debug("APNS notification successfully sent")
	return nil
}
