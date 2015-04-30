package slackboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"errors"
)

type SlackPayload struct {
	Channel   string `json:"channel"`
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	Text      string `json:"text"`
	Parse     string `json:"parse,omitempty"`
}

type SlackboardPayload struct {
	Tag  string `json:"tag"`
	Host string `json:"host,omitempty"`
	Text string `json:"text"`
	Sync bool   `json:"sync,omitempty"`
}

func sendNotification2Slack(payload *SlackPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Post(
		ConfSlackboard.Core.SlackURL,
		"application/json",
		strings.NewReader(string(body)))

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Slack is not available:%s", resp.Status))
	}

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func NotifyHandler(w http.ResponseWriter, r *http.Request) {
	LogError.Debug("notify-request is Accepted")

	LogError.Debug("parse request body")
	var req SlackboardPayload
	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		LogAcceptedRequest("/notify", r.Method, r.Proto, r.ContentLength, "")
		sendResponse(w, "Request-body is malformed", http.StatusBadRequest)
		return
	}

	LogAcceptedRequest("/notify", r.Method, r.Proto, r.ContentLength, req.Tag)

	LogError.Debug("method check")
	if r.Method != "POST" {
		sendResponse(w, "invalid method", http.StatusBadRequest)
		return
	}

	LogError.Debug("find tag")
	sent := false
	for i, tag := range ConfSlackboard.Tags {
		if tag.Tag == req.Tag {
			atomic.AddUint64(&Topics[i].Count, 1)
			payload := &SlackPayload{
				Channel:   tag.Channel,
				Username:  tag.Username,
				IconEmoji: tag.IconEmoji,
				Text:      req.Text,
				Parse:     tag.Parse,
			}

			if req.Sync {
				err := sendNotification2Slack(payload)
				if err != nil {
					sendResponse(w, "failed to post message to slack", http.StatusBadGateway)
					return
				}
				sent = true
			} else {
				go func() {
					err := sendNotification2Slack(payload)
					if err != nil {
						LogError.Error(fmt.Sprintf("failed to post message to slack:%s", err.Error()))
					}
				}()
			}
		}

	}

	LogError.Debug("response to client")

	if req.Sync {
		if sent {
			sendResponse(w, "ok", http.StatusOK)
		} else {
			msg := fmt.Sprintf("tag:%s is not found", req.Tag)
			sendResponse(w, msg, http.StatusBadRequest)
		}

	} else {
		sendResponse(w, "ok", http.StatusOK)
	}
}
