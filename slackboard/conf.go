package slackboard

import (
	"errors"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type ConfToml struct {
	Core SectionCore  `toml:"core"`
	Tags []SectionTag `toml:"tags"`
	Log  SectionLog   `toml:"log"`
}

type SectionCore struct {
	Port             string `toml:"port"`
	SlackURL         string `toml:"slack_url"`
	SlackToken       string `toml:"slack_token"`
	QPS              int    `toml:"qps"`
	MaxDelayDuration int    `toml:"max_delay_duration"`
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

func init() {
	ConfSlackboard = BuildDefaultConf()
}

func BuildDefaultConf() ConfToml {
	var conf ConfToml
	// Core
	conf.Core.Port = "29800"
	conf.Core.SlackURL = ""
	conf.Core.SlackToken = ""
	conf.Core.QPS = 0
	conf.Core.MaxDelayDuration = -1 // means an empty parameter
	// Log
	conf.Log.AccessLog = "stdout"
	conf.Log.ErrorLog = "stderr"
	conf.Log.Level = "error"
	return conf
}

func LoadConf(confPath string, confToml *ConfToml) error {
	doc, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	if err := toml.Unmarshal(doc, confToml); err != nil {
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
	}
	return nil
}
