package slackboard

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	statsGo "github.com/fukata/golang-stats-api-handler"
)

func init() {
	statsGo.PrettyPrintEnabled()
}

func RegisterAPIs() {
	http.HandleFunc("/notify", NotifyHandler)
	http.HandleFunc("/notify-directly", NotifyDirectlyHandler)
	http.HandleFunc("/app/config", ConfigAppHandler)
	http.HandleFunc("/stat/go", statsGo.Handler)
}

func Run() {
	// Listen TCP Port
	if _, err := strconv.Atoi(ConfSlackboard.Core.Port); err == nil {
		http.ListenAndServe(":"+ConfSlackboard.Core.Port, nil)
	}

	// Listen UNIX Socket
	if strings.HasPrefix(ConfSlackboard.Core.Port, "unix:/") {
		sockPath := ConfSlackboard.Core.Port[5:]
		fi, err := os.Lstat(sockPath)
		if err == nil && (fi.Mode()&os.ModeSocket) == os.ModeSocket {
			err := os.Remove(sockPath)
			if err != nil {
				log.Fatalf("failed to remove %s", sockPath)
			}
		}
		l, err := net.Listen("unix", sockPath)
		if err != nil {
			log.Fatalf("failed to listen: %s", sockPath)
		}
		http.Serve(l, nil)
	}

	log.Fatalf("port parameter is invalid: %s", ConfSlackboard.Core.Port)
}
