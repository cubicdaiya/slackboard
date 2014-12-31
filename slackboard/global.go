package slackboard

import (
	"github.com/Sirupsen/logrus"
	"html/template"
)

var (
	ConfSlackboard ConfToml
	LogAccess      *logrus.Logger
	LogError       *logrus.Logger
	IndexTemplate  *template.Template
)
