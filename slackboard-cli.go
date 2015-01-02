package main

import (
	"./slackboard"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func send(server string, payload *slackboard.SlackboardPayload) error {
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

func main() {

	version := flag.Bool("v", false, "slackboard version")
	server := flag.String("s", "", "slackboard server name")
	tag := flag.String("t", "", "slackboard tag name")
	flag.Parse()

	if *version {
		slackboard.PrintVersion()
		os.Exit(0)
	}

	if *server == "" {
		log.Fatal("Specify slackboard server name")
	}

	if *tag == "" {
		log.Fatal("Specify slackboard tag name")
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	var text bytes.Buffer
	io.Copy(&text, os.Stdin)
	payload := &slackboard.SlackboardPayload{
		Tag:  *tag,
		Host: hostname,
		Text: text.String(),
	}

	err = send(*server, payload)
	if err != nil {
		log.Fatal(err.Error())
	}
}
