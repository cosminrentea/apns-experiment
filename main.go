package main

import (
  apns "github.com/sideshow/apns2"
  "github.com/sideshow/apns2/certificate"
  "log"
)

func main() {
  cert, pemErr := certificate.FromP12File("development-certificate.p12", "")
  if pemErr != nil {
    log.Println("Certificate Error:", pemErr)
  }

  notification := &apns.Notification{}
  notification.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
  notification.Topic = "APNS Test"
  notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`)

  client := apns.NewClient(cert).Development()
  res, err := client.Push(notification)
  if err != nil {
    log.Println("Error:", err)
    return
  }
  log.Println("APNs ID:", res.ApnsID)
}