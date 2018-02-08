package slackboard

import (
	"html/template"

	"github.com/sirupsen/logrus"
)

var (
	ConfSlackboard ConfToml
	QPSEnd         *QPSPerSlackEndpoint
	LogAccess      *logrus.Logger
	LogError       *logrus.Logger
	IndexTemplate  *template.Template
)
