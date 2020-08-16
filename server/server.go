package server

import (
	"doorman/types"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger

func beepFile(w http.ResponseWriter, r *http.Request) {
	if validateReqFromTwilio(w, r) {
		http.ServeFile(w, r, "./resources/beep.mp3")
	} else {
		logger.Println("Non twilio request to beepFile")
	}
}

func knock(w http.ResponseWriter, r *http.Request) {
	if validateReqFromTwilio(w, r) {
		logger.Println("Someone knocked")
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
			logger.Println("Delivery person entered")
			twiml = types.TwiML{Say: os.Getenv("DELIVERY_WELCOME"), Play: &play}
			sendNotifications("Delivery person just entered the building")
		case os.Getenv("PERSONAL_PASSWORD"):
			logger.Println("Friend entered")
			twiml = types.TwiML{Say: os.Getenv("PERSONAL_WELCOME"), Play: &play}
			sendNotifications("A friend is here :)")
		default:
			logger.Println("Someone failed the password check with code: " + password)
			twiml = types.TwiML{Say: os.Getenv("DENIED_MESSAGE")}
			sendNotifications("Someone typed the wrong password for the building")
		}

		twilioWriter(twiml, w)
	} else {
		logger.Println("Non twilio request to attemptEntry")
	}
}

func StartServer() {
	f, err := os.OpenFile("/var/log/doorman.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger = log.New(f, "", log.LstdFlags)

	http.HandleFunc(os.Getenv("KNOCK"), knock)
	http.HandleFunc(os.Getenv("BEEP"), beepFile)
	http.HandleFunc(os.Getenv("ATTEMPT_ENTRY"), attemptEntry)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
