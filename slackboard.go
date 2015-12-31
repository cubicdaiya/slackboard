package main

import (
	"./slackboard"
	"flag"
	"log"
)

func main() {

	version := flag.Bool("v", false, "slackboard version")
	confPath := flag.String("c", "", "configuration file for slackboard")
	flag.Parse()

	if *version {
		slackboard.PrintVersion()
		return
	}

	// Load conf
	err := slackboard.LoadConf(*confPath, &slackboard.ConfSlackboard)
	if err != nil {
		log.Fatal(err)
	}

	// set logger
	err = slackboard.SetLogLevel(slackboard.LogAccess, "info")
	if err != nil {
		log.Fatal(err)
	}
	err = slackboard.SetLogLevel(slackboard.LogError, slackboard.ConfSlackboard.Log.Level)
	if err != nil {
		log.Fatal(err)
	}
	err = slackboard.SetLogOut(slackboard.LogAccess, slackboard.ConfSlackboard.Log.AccessLog)
	if err != nil {
		log.Fatal(err)
	}
	err = slackboard.SetLogOut(slackboard.LogError, slackboard.ConfSlackboard.Log.ErrorLog)
	if err != nil {
		log.Fatal(err)
	}

	slackboard.RegisterAPIs()
	slackboard.SetupUI()
	slackboard.Run()
}
