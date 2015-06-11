package main

import (
	"./slackboard"
	"flag"
	"github.com/Sirupsen/logrus"
	statsGo "github.com/fukata/golang-stats-api-handler"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func main() {

	version := flag.Bool("v", false, "slackboard version")
	confPath := flag.String("c", "", "configuration file for slackboard")
	flag.Parse()

	if *version {
		slackboard.PrintVersion()
		os.Exit(0)
	}

	// Set concurrency
	runtime.GOMAXPROCS(runtime.NumCPU())

	// init logger
	slackboard.LogAccess = logrus.New()
	slackboard.LogError = logrus.New()

	slackboard.LogAccess.Formatter = new(slackboard.SlackboardFormatter)
	slackboard.LogError.Formatter = new(slackboard.SlackboardFormatter)

	// Load conf
	slackboard.ConfSlackboard = slackboard.BuildDefaultConf()
	err := slackboard.LoadConf(*confPath, &slackboard.ConfSlackboard)
	if err != nil {
		log.Fatal(err.Error())
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

	statsGo.PrettyPrintEnabled()
	http.HandleFunc("/notify", slackboard.NotifyHandler)
	http.HandleFunc("/notify-directly", slackboard.NotifyDirectlyHandler)
	http.HandleFunc("/app/config", slackboard.ConfigAppHandler)
	http.HandleFunc("/stat/go", statsGo.Handler)
	slackboard.SetupUI()

	// Listen TCP Port
	if _, err := strconv.Atoi(slackboard.ConfSlackboard.Core.Port); err == nil {
		http.ListenAndServe(":"+slackboard.ConfSlackboard.Core.Port, nil)
	}

	// Listen UNIX Socket
	if strings.HasPrefix(slackboard.ConfSlackboard.Core.Port, "unix:/") {
		sockPath := slackboard.ConfSlackboard.Core.Port[5:]
		fi, err := os.Lstat(sockPath)
		if err == nil && (fi.Mode()&os.ModeSocket) == os.ModeSocket {
			err := os.Remove(sockPath)
			if err != nil {
				log.Fatal("failed to remove " + sockPath)
			}
		}
		l, err := net.Listen("unix", sockPath)
		if err != nil {
			log.Fatal("failed to listen: " + sockPath)
		}
		http.Serve(l, nil)
	}

	log.Fatal("port parameter is invalid: " + slackboard.ConfSlackboard.Core.Port)

}
