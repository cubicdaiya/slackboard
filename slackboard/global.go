package slackboard

import (
	"github.com/sirupsen/logrus"
)

var (
	ConfSlackboard ConfToml
	QPSEnd         *QPSPerSlackEndpoint
	LogAccess      *logrus.Logger
	LogError       *logrus.Logger
)
