package slackboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func sendNotification2Slackboard(server, api, body string) error {
	client := &http.Client{}

	url := fmt.Sprintf("http://%s/%s", server, api)
	resp, err := client.Post(
		url,
		"application/json",
		strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(content))
	}

	return nil
}

func SendNotification2SlackboardDirectly(server string, payload *SlackboardDirectPayload) error {
	if strings.Index(payload.Payload.Channel, "#") != 0 {
		payload.Payload.Channel = "#" + payload.Payload.Channel
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return sendNotification2Slackboard(server, "notify-directly", string(body))
}

func SendNotification2Slackboard(server string, payload *SlackboardPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return sendNotification2Slackboard(server, "notify", string(body))
}
