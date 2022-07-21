package main

import (
	"logger-service/data"
	"net/http"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload jsonPayload
	app.readJson(w, r, &requestPayload)

	//insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJson(w, err, 400)
		return
	}
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
		Data:    event,
	}
	app.writeJson(w, http.StatusAccepted, resp)
}
