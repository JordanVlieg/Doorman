package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"frontdoor/types"

	_ "github.com/joho/godotenv/autoload"
	twilio_go "github.com/kevinburke/twilio-go"
)

var logger *log.Logger

func beepFile(w http.ResponseWriter, r *http.Request) {
	if validateReqFromTwilio(w, r) {
		http.ServeFile(w, r, "./beep.mp3")
	} else {
		logger.Println("Non twilio request to beepFile")
	}
}

func knock(w http.ResponseWriter, r *http.Request) {
	if validateReqFromTwilio(w, r) {
		play := types.Play{Message: getFullURI("BEEP")}
		gather := types.Gather{Action: os.Getenv("ATTEMPT_ENTRY"), Method: "POST", NumDigits: "4", Timeout: "7", ActionOnEmptyResult: "false", Say: os.Getenv("MAIN_WELCOME"), Play: &play}
		twiml := types.TwiML{Gather: &gather}
		twilioWriter(twiml, w)
	} else {
		logger.Println("Non twilio request to knock")
	}
}

func attemptEntry(w http.ResponseWriter, r *http.Request) {
	if validateReqFromTwilio(w, r) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Println(err.Error())
			return
		}

		play := types.Play{Digits: doorCodeToTwilioTone(os.Getenv("BUZZ_CODE"))}
		var twiml types.TwiML

		switch password := r.Form["Digits"][0]; password {
		case os.Getenv("DELIVERY_PASSWORD"):
			twiml = types.TwiML{Say: os.Getenv("DELIVERY_WELCOME"), Play: &play}
		case os.Getenv("PERSONAL_PASSWORD"):
			twiml = types.TwiML{Say: os.Getenv("PERSONAL_WELCOME"), Play: &play}
		default:
			twiml = types.TwiML{Say: os.Getenv("DENIED_MESSAGE")}
		}

		twilioWriter(twiml, w)
	} else {
		logger.Println("Non twilio request to attemptEntry")
	}
}

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
		err := twilio_go.ValidateIncomingRequest(os.Getenv("BASE_URI"), os.Getenv("AUTH_TOKEN"), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			logger.Println(err.Error())
			return false
		}
	}
	return true
}

func main() {
	f, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger = log.New(f, "prefix", log.LstdFlags)

	http.HandleFunc(os.Getenv("KNOCK"), knock)
	http.HandleFunc(os.Getenv("BEEP"), beepFile)
	http.HandleFunc(os.Getenv("ATTEMPT_ENTRY"), attemptEntry)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
