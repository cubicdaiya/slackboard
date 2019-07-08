package slackboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

type UnixDial struct {
	sockPath string
}

type Retry struct {
	WaitMin time.Duration
	WaitMax time.Duration
	Max     int
}

func (u UnixDial) Dial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", u.sockPath)
}

func sendNotification2Slackboard(server, api, body string, retry *Retry) error {
	var client *retryablehttp.Client
	var url string

	if strings.HasPrefix(server, "unix:/") {
		// UNIX Socket
		client = &retryablehttp.Client{
			HTTPClient: &http.Client{
				Transport: &http.Transport{
					Dial: (&UnixDial{sockPath: server[5:]}).Dial,
				},
				Timeout: 30 * time.Second,
			},
			RetryWaitMin: retry.WaitMin,
			RetryWaitMax: retry.WaitMax,
			RetryMax:     retry.Max,
		}
		url = fmt.Sprintf("http://localhost/%s", api)
	} else {
		// TCP
		client = &retryablehttp.Client{
			HTTPClient: &http.Client{
				Timeout: 30 * time.Second,
			},
			RetryWaitMin: retry.WaitMin,
			RetryWaitMax: retry.WaitMax,
			RetryMax:     retry.Max,
		}

		url = fmt.Sprintf("http://%s/%s", server, api)
		if strings.HasPrefix(server, "http://") || strings.HasPrefix(server, "https://") {
			url = fmt.Sprintf("%s/%s", server, api)
		}
	}

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

func SendNotification2SlackboardDirectly(server string, payload *SlackboardDirectPayload, retry *Retry) error {
	if strings.Index(payload.Payload.Channel, "#") != 0 &&
		strings.Index(payload.Payload.Channel, "@") != 0 {
		payload.Payload.Channel = "#" + payload.Payload.Channel
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return sendNotification2Slackboard(server, "notify-directly", string(body), retry)
}

func SendNotification2Slackboard(server string, payload *SlackboardPayload, retry *Retry) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return sendNotification2Slackboard(server, "notify", string(body), retry)
}
