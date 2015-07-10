package main

import (
	"./slackboard"
	"bytes"
	"flag"
	"io"
	"log"
	"os"
)

func main() {

	version := flag.Bool("v", false, "slackboard version")
	server := flag.String("s", "", "slackboard server name")
	tag := flag.String("t", "", "slackboard tag name")
	sync := flag.Bool("sync", false, "enable synchronous notification")
	channel := flag.String("c", "", "slackboard channel name")
	username := flag.String("u", "slackboard", "user name")
	iconemoji := flag.String("i", ":clipboard:", "emoji icon")
	parse := flag.String("p", "full", "parsing mode")
	flag.Parse()

	if *version {
		slackboard.PrintVersion()
		os.Exit(0)
	}

	if *server == "" && *tag == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *server == "" {
		log.Fatal("Specify slackboard server name")
	}

	if *tag == "" && *channel == "" {
		log.Fatal("Specify slackboard tag or channel name")
	}

	if *tag != "" && *channel != "" {
		log.Fatal("Assigning with '-t' and '-c' at once is not allowed")
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	var text bytes.Buffer
	io.Copy(&text, os.Stdin)

	if *tag != "" {

		payload := &slackboard.SlackboardPayload{
			Tag:  *tag,
			Host: hostname,
			Text: text.String(),
			Sync: *sync,
		}

		err = slackboard.SendNotification2Slackboard(*server, payload)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	payloadSlack := slackboard.SlackPayload{
		Channel:   *channel,
		Username:  *username,
		IconEmoji: *iconemoji,
		Text:      text.String(),
		Parse:     *parse,
	}
	payloadDirectly := &slackboard.SlackboardDirectPayload{
		Payload: payloadSlack,
		Sync:    *sync,
	}

	err = slackboard.SendNotification2SlackboardDirectly(*server, payloadDirectly)
	if err != nil {
		log.Fatal(err)
	}

}
