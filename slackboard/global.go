package slackboard

import (
	"github.com/Sirupsen/logrus"
	"html/template"
)

var (
	ConfSlackboard ConfToml
	Topics         []Topic
	LogAccess      *logrus.Logger
	LogError       *logrus.Logger
	IndexTemplate  *template.Template
)
