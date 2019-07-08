package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/cubicdaiya/slackboard/slackboard"
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
	level := flag.String("l", "message", "level")
	color := flag.String("C", "", "color")
	title := flag.String("title", "", "title")
	retryMax := flag.Int("retry-max", 4, "Maximum number of retries")
	retryWaitMin := flag.Duration("retry-wait-min", 1*time.Second, "Minimum time to wait when retry")
	retryWaitMax := flag.Duration("retry-wait-max", 10*time.Second, "Maximum time to wait when retry")
	flag.Parse()

	if *version {
		slackboard.PrintVersion()
		return
	}

	if *server == "" {
		*server = os.Getenv("SLACKBOARD_SERVER")
	}

	if *server == "" && *tag == "" {
		flag.PrintDefaults()
		return
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

	retry := &slackboard.Retry{
		WaitMin: *retryWaitMin,
		WaitMax: *retryWaitMax,
		Max:     *retryMax,
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	var text bytes.Buffer
	io.Copy(&text, os.Stdin)

	if *tag != "" {

		payload := &slackboard.SlackboardPayload{
			Tag:   *tag,
			Host:  hostname,
			Text:  text.String(),
			Sync:  *sync,
			Level: *level,
			Title: *title,
		}

		err = slackboard.SendNotification2Slackboard(*server, payload, retry)
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

	if *color == "" && *level != "" {
		switch *level {
		case "info":
			*color = "#00ff00"
		case "warn":
			*color = "#ffdd00"
		case "crit":
			*color = "#ff0000"
		}
	}

	if *color != "" || *title != "" {
		payloadSlack.Text = ""
		payloadSlack.Attachments = make([]slackboard.SlackPayloadAttachments, 1)
		payloadSlack.Attachments[0] = slackboard.SlackPayloadAttachments{
			Color: *color,
			Title: *title,
			Text:  text.String(),
		}
	}

	payloadDirectly := &slackboard.SlackboardDirectPayload{
		Payload: payloadSlack,
		Sync:    *sync,
	}

	err = slackboard.SendNotification2SlackboardDirectly(*server, payloadDirectly, retry)
	if err != nil {
		log.Fatal(err)
	}

}
