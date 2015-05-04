package slackboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendNotification2Slackboard(server string, payload *SlackboardPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}

	url := fmt.Sprintf("http://%s/notify", server)
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
