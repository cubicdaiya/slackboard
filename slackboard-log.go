package main

import (
	"./slackboard"
	"flag"
	"fmt"
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
	notify := flag.Bool("notity", true, "enable notification to slackboard")
	logfile := flag.String("log", "", "log-file path")
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

	if *tag == "" {
		log.Fatal("Specify slackboard tag name")
	}

	argv := flag.Args()

	if len(argv) == 0 {
		log.Fatal("command is not specified")
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	out, err := exec.Command(argv[0], argv[1:]...).CombinedOutput()
	if err != nil && *notify {
		text := fmt.Sprintf(`
Host   : %s
Command: %s
Output : %s
Error  : %s
`, hostname, strings.Join(argv, " "), string(out), err.Error())

		payload := &slackboard.SlackboardPayload{
			Tag:  *tag,
			Host: hostname,
			Text: text,
			Sync: *sync,
		}

		err = slackboard.SendNotification2Slackboard(*server, payload)
		if err != nil {
			log.Fatal(err.Error())
		}
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
			log.Fatal(err.Error())
		}
		defer file.Close()
		file.WriteString(string(out))
	} else {
		file, err := os.Create(*logfile)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()
		file.WriteString(string(out))
	}

}
