package cluster

import "github.com/sirupsen/logrus"

var (
	ctxLogger  = logrus.WithField("m", "context")
	connLogger = logrus.WithField("m", "conn")
	nodeLogger = logrus.WithField("m", "node")
)
