package slackboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"time"
)

type LogReq struct {
	Type          string `json:"type"`
	Time          string `json:"time"`
	URI           string `json:"uri"`
	Method        string `json:"method"`
	Proto         string `json:"proto"`
	ContentLength int64  `json:"content_length"`
}

type SlackboardFormatter struct {
}

func (f *SlackboardFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "[%s] ", entry.Level.String())
	fmt.Fprintf(b, "%s", entry.Message)
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func InitLog() *logrus.Logger {
	return logrus.New()
}

func SetLogOut(log *logrus.Logger, outString string) error {
	switch outString {
	case "stdout":
		log.Out = os.Stdout
	case "stderr":
		log.Out = os.Stderr
	default:
		f, err := os.OpenFile(outString, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		log.Out = f
	}
	return nil
}

func SetLogLevel(log *logrus.Logger, levelString string) error {
	level, err := logrus.ParseLevel(levelString)
	if err != nil {
		return err
	}
	log.Level = level
	return nil
}

func LogAcceptedRequest(uri, method, proto string, length int64) {
	log := &LogReq{
		Type:          "accepted-request",
		Time:          time.Now().Format("2006/01/02 15:04:05 MST"),
		URI:           uri,
		Method:        method,
		Proto:         proto,
		ContentLength: length,
	}
	logJSON, err := json.Marshal(log)
	if err != nil {
		LogError.Error("Marshaling JSON error")
		return
	}
	LogAccess.Info(string(logJSON))
}
