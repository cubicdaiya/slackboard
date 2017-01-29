package slackboard

import (
	"html/template"

	"github.com/Sirupsen/logrus"
)

var (
	ConfSlackboard ConfToml
	QPSEnd         *QPSPerSlackEndpoint
	Topics         []Topic
	LogAccess      *logrus.Logger
	LogError       *logrus.Logger
	IndexTemplate  *template.Template
)
