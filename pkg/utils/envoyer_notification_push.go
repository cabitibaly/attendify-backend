package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type ExpoMessage struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var url string

func InitializePushURL(urlEnv string) {
	url = urlEnv
}

func EnvoyerNotificationPush(token, title, body string) error {
	message := ExpoMessage{
		To:    token,
		Title: title,
		Body:  body,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("une erreur est survenue lors de l'envoi de la notification")
	}

	defer response.Body.Close()

	return nil

}
