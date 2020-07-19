package main

import (
	"encoding/xml"
	"net/http"
	"os"

	"frontdoor/types"

	_ "github.com/joho/godotenv/autoload"
	twilio_go "github.com/kevinburke/twilio-go"
)

func beepFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./beep.mp3")
}

func knock(w http.ResponseWriter, r *http.Request) {
	play := types.Play{Message: getFullURI("BEEP")}
	gather := types.Gather{Action: os.Getenv("ATTEMPT_ENTRY"), Method: "POST", NumDigits: "4", Timeout: "7", ActionOnEmptyResult: "false", Say: os.Getenv("MAIN_WELCOME"), Play: &play}
	twiml := types.TwiML{Gather: &gather}
	twilioWriter(twiml, w)
}

func attemptEntry(w http.ResponseWriter, r *http.Request) {
	err := twilio_go.ValidateIncomingRequest(os.Getenv("BASE_URI"), os.Getenv("AUTH_TOKEN"), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	play := types.Play{Digits: "9w9w"}
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
}

func getFullURI(pathEnvVar string) string {
	return os.Getenv("BASE_URI") + os.Getenv(pathEnvVar)
}

func twilioWriter(twiml types.TwiML, w http.ResponseWriter) {
	x, err := xml.Marshal(twiml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

func main() {
	http.HandleFunc(os.Getenv("KNOCK"), knock)
	http.HandleFunc(os.Getenv("BEEP"), beepFile)
	http.HandleFunc(os.Getenv("ATTEMPT_ENTRY"), attemptEntry)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
