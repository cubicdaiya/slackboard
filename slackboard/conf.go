package slackboard

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/BurntSushi/toml"
)

type ConfToml struct {
	Core SectionCore  `toml:"core"`
	Tags []SectionTag `toml:"tags"`
	Log  SectionLog   `toml:"log"`
	UI   SectionUI    `toml:"ui"`
}

type SectionCore struct {
	Port     string `toml:"port"`
	SlackURL string `toml:"slack_url"`
}

type SectionTag struct {
	Tag       string `toml:"tag"`
	Channel   string `toml:"channel"`
	Username  string `toml:"username"`
	IconEmoji string `toml:"icon_emoji"`
	Parse     string `toml:"parse"`
}

type SectionLog struct {
	AccessLog string `toml:"access_log"`
	ErrorLog  string `toml:"error_log"`
	Level     string `toml:"level"`
}

type SectionUI struct {
	Root string `toml:"root"`
}

func init() {
	ConfSlackboard = BuildDefaultConf()
}

func BuildDefaultConf() ConfToml {
	var conf ConfToml
	// Core
	conf.Core.Port = "29800"
	conf.Core.SlackURL = ""
	// Log
	conf.Log.AccessLog = "stdout"
	conf.Log.ErrorLog = "stderr"
	conf.Log.Level = "error"
	return conf
}

func LoadConf(confPath string, confToml *ConfToml) error {
	_, err := toml.DecodeFile(confPath, confToml)
	if err != nil {
		return err
	}
	for i, tag := range confToml.Tags {
		if tag.Tag == "" {
			return errors.New("tag is empty")
		}
		if tag.Channel == "" {
			confToml.Tags[i].Channel = "#random"
		}
		if tag.Username == "" {
			confToml.Tags[i].Username = "slackboard"
		}
		if tag.IconEmoji == "" {
			confToml.Tags[i].IconEmoji = ":clipboard:"
		}
		if tag.Parse == "" {
			confToml.Tags[i].Parse = "full"
		}
		topic := Topic{Tag: confToml.Tags[i], Count: 0}
		Topics = append(Topics, topic)
	}
	return nil
}

func ConfigAppHandler(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	e := toml.NewEncoder(&b)
	result := ConfSlackboard
	result.Core.SlackURL = "..."
	err := e.Encode(result)
	if err != nil {
		msg := "Response-body could not be created"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, b.String())
}
