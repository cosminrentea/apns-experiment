package main

import (
	"fmt"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/smancke/guble/protocol"
	"github.com/smancke/guble/server"
	"log"
)

func init() {
	server.AfterMessageDelivery = func(m *protocol.Message) {
		fmt.Print("message delivered")
	}
}

func main() {

	// server.Main()

	cert, errCert := certificate.FromP12File("development-certificate.p12", "WeLoveApple")
	if errCert != nil {
		log.Println("Certificate Error: ", errCert)
	}

	notification := &apns.Notification{}
	notification.DeviceToken = "dgELOSlqfW0:APA91bHxaLpeQzqKyDecIWKLahLhe_H2vPCqIxpqAqOR7FQWTV-QeNRuPCtLNFnrlwMTiAWGyhwQji5G5FuqvQ0V7qPgDSaTBybdSJdg21ss2613tflHLA3QJWuBDNU1n9KmpqixfhOV"
	notification.Topic = "com.rewe.iosapp"
	notification.Payload = []byte(`{"aps":{"alert":"Hello from guble team, using APNS!"}}`)

	client := apns.NewClient(cert).Development()
	res, errPush := client.Push(notification)
	if errPush != nil {
		log.Println("Error:", errPush)
		return
	}
	log.Println("APNS ID:", res.ApnsID)
}
