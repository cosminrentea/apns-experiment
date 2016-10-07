package main

import (
	"errors"
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
	msgAPSNNotSent      = "APNS notification was not sent"
	errAPSNNotSent      = errors.New(msgAPSNNotSent)
)

func init() {
	server.AfterMessageDelivery = func(m *protocol.Message) {
		fmt.Print("message delivered")
	}
}

type APNSConfig struct {
	CertFileName string
	CertPassword string
}

func main() {

	// server.Main()

	cfg := APNSConfig{
		CertFileName: defaultCertFileName,
		CertPassword: os.Getenv("APNS_CERT_PASSWORD"),
	}
	topic := os.Getenv("APNS_TOPIC")
	deviceToken := os.Getenv("APNS_DEVICE_TOKEN")
	cl := getAPNSClient(cfg, false)
	p := payload.NewPayload().
		AlertTitle("REWE Sonderrabatt").
		AlertBody("Sie haben ein Sonderrabatt von 50% für das neue iPhone 8 bekommen!").
		ContentAvailable()
	sendAPNSNotification(cl, topic, deviceToken, p)
}

func getAPNSClient(cfg *APNSConfig, production bool) *apns.Client {
	cert, errCert := certificate.FromP12File(cfg.CertFileName, cfg.CertPassword)
	if errCert != nil {
		log.WithError(errCert).Error("APNS certificate error")
	}
	if production {
		return apns.NewClient(cert).Production()
	}
	return apns.NewClient(cert).Development()
}

func sendAPNSNotification(cl *apns.Client, topic string, deviceToken string, p *payload.Payload) error {
	notification := &apns.Notification{
		Priority:    apns.PriorityHigh,
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
		log.WithField("id", response.ApnsID).Error(msgAPSNNotSent)
		return errAPSNNotSent
	}
	log.WithField("id", response.ApnsID).Debug("APNS notification successfully sent")
	return nil
}
