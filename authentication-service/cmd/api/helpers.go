package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Data    any    `jaons:"data"`
	Message string `json:"message"`
}

func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxByte := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByte))
	// log.Printf("\n\nRequest Body :\n %s\n\n", r.Body)
	dec := json.NewDecoder(r.Body)
	// log.Printf("\n\nDecoded Body :\n %s\n\n", dec)

	err := dec.Decode(data)
	log.Printf("\n\n data passed in Body :\n %s\n\n", data)

	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must have only single json value")
	}
	return nil
}

func (app *Config) writeJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	log.Printf("🐼: %s", data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()
	log.Printf("\nError: %s \n", payload)

	return app.writeJson(w, statusCode, payload)
}
