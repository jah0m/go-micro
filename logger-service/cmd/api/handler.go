package main

import (
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var payload JSONPayload
	_ = app.readJSON(w, r, &payload)

	//insert the log into the database
	err := app.Models.LogEntry.Insert(data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	})

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "log entry created",
	}

	app.writeJSON(w, http.StatusCreated, resp)

}
