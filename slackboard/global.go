package slackboard

import (
	"html/template"

	"github.com/Sirupsen/logrus"
)

var (
	ConfSlackboard ConfToml
	Topics         []Topic
	LogAccess      *logrus.Logger
	LogError       *logrus.Logger
	IndexTemplate  *template.Template
)
