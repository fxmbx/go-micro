package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
		Data:    "Palazzo, jiggy, gb'oja, k'ole reason",
	}

	_ = app.writeJson(w, http.StatusOK, payload)

	// out, _ := json.MarshalIndent(payload, "", "\t")
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusAccepted)
	// w.Write(out)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	log.Printf("request Payload 😁 %v", requestPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendmail(w, requestPayload.Mail)

	default:
		app.errorJson(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	log.Printf("auth Payload 😁 %s", jsonData)

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
	}

	defer response.Body.Close()
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Println(response.StatusCode)
		app.errorJson(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJson(w, errors.New(jsonFromService.Message))
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}
	app.writeJson(w, http.StatusAccepted, payload)

}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		app.errorJson(w, err)
		// log.Println(err)
		// return
	}

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return

	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("something went wrong: 😞 "))
	}
	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJson(w, errors.New(jsonFromService.Message))
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged"
	payload.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, payload)

}

func (app *Config) sendmail(w http.ResponseWriter, msg MailPayload) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	request, err := http.NewRequest("POST", "http://mailer-service/send", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("something went wrong: 😞 "))
		return
	}
	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	app.errorJson(w, errors.New("something went wrong: 😞 "))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	if jsonFromService.Error {
		app.errorJson(w, errors.New(jsonFromService.Message))
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Mail Sent to " + msg.To

	app.writeJson(w, http.StatusAccepted, payload)

}
