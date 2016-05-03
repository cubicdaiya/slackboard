package slackboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
)

type SlackPayloadAttachmentsField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type SlackPayloadAttachments struct {
	Fallback string `json:"fallback"`
	Color    string `json:"color"`
	Pretext  string `json:"pretext"`

	AuthorName string `json:"author_name"`
	AuthorLink string `json:"author_link"`
	AuthorIcon string `json:"author_icon"`

	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`

	Field []SlackPayloadAttachmentsField `json:"fields"`

	ImageUrl string `json:"image_url"`
	ThumbUrl string `json:"thumb_url"`
}

type SlackPayload struct {
	Channel     string                    `json:"channel"`
	Username    string                    `json:"username,omitempty"`
	IconEmoji   string                    `json:"icon_emoji,omitempty"`
	Text        string                    `json:"text"`
	Parse       string                    `json:"parse,omitempty"`
	Attachments []SlackPayloadAttachments `json:"attachments"`
}

type SlackboardPayload struct {
	Tag   string `json:"tag"`
	Host  string `json:"host,omitempty"`
	Text  string `json:"text"`
	Sync  bool   `json:"sync,omitempty"`
	Level string `json:"level"`
	Title string `json:"title,omitempty"`
}

type SlackboardDirectPayload struct {
	Payload SlackPayload `json:"payload"`
	Sync    bool         `json:"sync,omitempty"`
}

func sendNotification2Slack(payload *SlackPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := http.DefaultClient

	resp, err := client.Post(
		ConfSlackboard.Core.SlackURL,
		"application/json",
		strings.NewReader(string(body)))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack is not available:%s", resp.Status)
	}

	return nil
}

func NotifyHandler(w http.ResponseWriter, r *http.Request) {
	LogError.Debug("notify-request is Accepted")

	LogError.Debug("parse request body")
	var req SlackboardPayload
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogAcceptedRequest(r, "")
		sendResponse(w, "failed to read request-body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		LogAcceptedRequest(r, "")
		sendResponse(w, "Request-body is malformed", http.StatusBadRequest)
		return
	}

	LogAcceptedRequest(r, req.Tag)

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

			var (
				color string
			)

			levelToColorMap := map[string]string{
				"info": "#00ff00", // green
				"warn": "#ffdd00", // yellow
				"crit": "#ff0000", // red
			}

			if color_, ok := levelToColorMap[req.Level]; ok {
				payload.Text = ""
				color = color_
			}

			if req.Title != "" {
				payload.Text = ""
			}

			if color != "" || req.Title != "" {
				payload.Attachments = make([]SlackPayloadAttachments, 1)
				payload.Attachments[0] = SlackPayloadAttachments{
					Color: color,
					Title: req.Title,
					Text:  req.Text,
				}
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

func NotifyDirectlyHandler(w http.ResponseWriter, r *http.Request) {
	LogError.Debug("notify-directly-request is Accepted")

	LogError.Debug("parse request body")
	var req SlackboardDirectPayload
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogAcceptedRequest(r, "")
		sendResponse(w, "failed to read request-body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		LogAcceptedRequest(r, "")
		sendResponse(w, "Request-body is malformed", http.StatusBadRequest)
		return
	}

	LogAcceptedRequest(r, req.Payload.Channel)

	LogError.Debug("method check")
	if r.Method != "POST" {
		sendResponse(w, "invalid method", http.StatusBadRequest)
		return
	}

	if req.Sync {
		err := sendNotification2Slack(&req.Payload)
		if err != nil {
			sendResponse(w, "failed to post message to slack", http.StatusBadGateway)
			return
		}
	} else {
		go func() {
			err := sendNotification2Slack(&req.Payload)
			if err != nil {
				LogError.Error(fmt.Sprintf("failed to post message to slack:%s", err.Error()))
			}
		}()
	}

	LogError.Debug("response to client")
	sendResponse(w, "ok", http.StatusOK)
}
