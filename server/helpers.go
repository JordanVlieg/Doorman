package server

import (
	"doorman/types"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"

	twilio_go "github.com/kevinburke/twilio-go"
)

func doorCodeToTwilioTone(code string) string {
	var builder strings.Builder
	builder.Grow(len(code) * 4)
	for _, c := range code {
		fmt.Fprintf(&builder, "%c%c%cwww", c, c, c)
	}
	return builder.String()
}

func getFullURI(pathEnvVar string) string {
	return os.Getenv("BASE_URI") + os.Getenv(pathEnvVar)
}

func sendNotifications(message string) {
	client := twilio_go.NewClient(os.Getenv("TWILIO_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), nil)
	// Send a message
	if os.Getenv("NOTIFY_PHONE_NUMBER") != "" {
		_, err := client.Messages.SendMessage(os.Getenv("TWILIO_PHONE_NUMBER"), os.Getenv("NOTIFY_PHONE_NUMBER"), message, nil)
		if err != nil {
			logger.Println("Failed to send notification SMS")
			return
		}
		logger.Println("Notified primary phone number: '" + message + "'")
	}
	if os.Getenv("SECONDARY_NOTIFY_NUMBER") != "" {
		_, err := client.Messages.SendMessage(os.Getenv("TWILIO_PHONE_NUMBER"), os.Getenv("SECONDARY_NOTIFY_NUMBER"), message, nil)
		if err != nil {
			logger.Println("Failed to send secondary notification SMS")
			return
		}
		logger.Println("Notified secondary phone number: '" + message + "'")
	}
}

func twilioWriter(twiml types.TwiML, w http.ResponseWriter) {
	x, err := xml.Marshal(twiml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

func validateReqFromTwilio(w http.ResponseWriter, r *http.Request) bool {
	if os.Getenv("DEVELOPMENT") != "true" {
		err := twilio_go.ValidateIncomingRequest(os.Getenv("BASE_URI"), os.Getenv("TWILIO_AUTH_TOKEN"), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			logger.Println(err.Error())
			return false
		}
	}
	return true
}
