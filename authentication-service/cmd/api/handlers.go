package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("handling authentication")
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJson(w, r, &requestPayload)
	log.Printf("\n%s\n", requestPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJson(w, errors.New("inavlid credentials 📧"))
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJson(w, errors.New("inavlid credentials 🔑"))
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = fmt.Sprintf("logged in user %s", user.Email)
	payload.Data = user

	app.writeJson(w, http.StatusAccepted, payload)
}
