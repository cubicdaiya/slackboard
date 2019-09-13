package slackboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	SlackPostMessageAPIURL = `https://slack.com/api/chat.postMessage`
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

func sendNotification2Slack(payload *SlackPayload, sync bool) (int, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return http.StatusBadGateway, err
	}

	if QPSEnd != nil {
		if sync && !QPSEnd.Available() {
			LogError.Warnf("Reject a sync message due to ratelimit:%s", string(body))
			return http.StatusTooManyRequests, fmt.Errorf("QPS ratelimit error")
		}
		if !sync && !QPSEnd.WaitAndAvailable() {
			LogError.Warnf("Reject an async message due to ratelimit:%s", string(body))
			return http.StatusTooManyRequests, fmt.Errorf("QPS ratelimit error")
		}
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if len(ConfSlackboard.Core.SlackToken) == 0 {
		// if SlackToken is not specified, use Incoming Webhook
		resp, err := client.Post(
			ConfSlackboard.Core.SlackURL,
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			return http.StatusBadGateway, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return http.StatusBadGateway, fmt.Errorf("Slack is not available:%s", resp.Status)
		}

		return http.StatusOK, nil
	}

	req, err := http.NewRequest(
		http.MethodPost,
		SlackPostMessageAPIURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return http.StatusBadGateway, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ConfSlackboard.Core.SlackToken)

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusBadGateway, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return http.StatusBadGateway, fmt.Errorf("Slack is not available:%s", resp.Status)
	}

	var errResp struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return http.StatusBadGateway, fmt.Errorf("Slack returned invalid format response, should be json: %s", err)
	}
	if !errResp.OK {
		return http.StatusBadGateway, fmt.Errorf("Slack returned error response: %s", errResp.Error)
	}

	return http.StatusOK, nil
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
	if r.Method != http.MethodPost {
		sendResponse(w, "invalid method", http.StatusBadRequest)
		return
	}

	LogError.Debug("find tag")
	sent := false
	for _, tag := range ConfSlackboard.Tags {
		if tag.Tag == req.Tag {
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
				status, err := sendNotification2Slack(payload, req.Sync)
				if err != nil {
					sendResponse(w, "failed to post message to slack", status)
					return
				}
				sent = true
			} else {
				go func() {
					_, err := sendNotification2Slack(payload, req.Sync)
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
	if r.Method != http.MethodPost {
		sendResponse(w, "invalid method", http.StatusBadRequest)
		return
	}

	if req.Sync {
		status, err := sendNotification2Slack(&req.Payload, req.Sync)
		if err != nil {
			sendResponse(w, "failed to post message to slack", status)
			return
		}
	} else {
		go func() {
			_, err := sendNotification2Slack(&req.Payload, req.Sync)
			if err != nil {
				LogError.Error(fmt.Sprintf("failed to post message to slack:%s", err.Error()))
			}
		}()
	}

	LogError.Debug("response to client")
	sendResponse(w, "ok", http.StatusOK)
}
