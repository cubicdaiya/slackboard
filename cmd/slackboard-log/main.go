package main

import (
	"flag"
	"fmt"
	"github.com/cubicdaiya/slackboard/slackboard"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	version := flag.Bool("v", false, "slackboard version")
	server := flag.String("s", "", "slackboard server name")
	tag := flag.String("t", "", "slackboard tag name")
	sync := flag.Bool("sync", false, "enable synchronous notification")
	notify := flag.Bool("notify", true, "enable notification to slackboard")
	logfile := flag.String("log", "", "log-file path")
	channel := flag.String("c", "", "slackboard channel name")
	username := flag.String("u", "slackboard", "user name")
	iconemoji := flag.String("i", ":clipboard:", "emoji icon")
	parse := flag.String("p", "full", "parsing mode")
	flag.Parse()

	if *version {
		slackboard.PrintVersion()
		return
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

	argv := flag.Args()

	if len(argv) == 0 {
		log.Fatal("command is not specified")
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	msgFmt := `
Host   : %s
Command: %s
Output : %s
Error  : %s
`

	out, err := exec.Command(argv[0], argv[1:]...).CombinedOutput()
	if err != nil && *notify {
		text := fmt.Sprintf(msgFmt,
			hostname,
			strings.Join(argv, " "),
			strings.TrimRight(string(out), "\n"),
			strings.TrimRight(err.Error(), "\n"))

		log.Println(text)

		if *tag != "" {
			payload := &slackboard.SlackboardPayload{
				Tag:   *tag,
				Host:  hostname,
				Text:  text,
				Sync:  *sync,
				Level: "crit",
			}

			err = slackboard.SendNotification2Slackboard(*server, payload)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			payloadSlack := slackboard.SlackPayload{
				Channel:   *channel,
				Username:  *username,
				IconEmoji: *iconemoji,
				Text:      text,
				Parse:     *parse,
			}
			payloadSlack.Text = ""
			payloadSlack.Attachments = make([]slackboard.SlackPayloadAttachments, 1)
			payloadSlack.Attachments[0] = slackboard.SlackPayloadAttachments{
				Color: "#ff0000",
				Text:  text,
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
	} else if err != nil {
		text := fmt.Sprintf(msgFmt, hostname, strings.Join(argv, " "), strings.TrimRight(string(out), "\n"), strings.TrimRight(err.Error(), "\n"))
		log.Println(text)
	}

	if *logfile == "" {
		return
	}

	fi, err := os.Stat(*logfile)
	if err == nil {
		if fi.IsDir() {
			log.Fatalf("%s is a directory.", *logfile)
		}
		file, err := os.OpenFile(*logfile, os.O_RDWR|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString(string(out))
	} else {
		file, err := os.Create(*logfile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString(string(out))
	}

}
