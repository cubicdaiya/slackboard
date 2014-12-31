package slackboard

import (
	"github.com/Sirupsen/logrus"
)

var (
	ConfSlackboard ConfToml
	LogAccess      *logrus.Logger
	LogError       *logrus.Logger
)
