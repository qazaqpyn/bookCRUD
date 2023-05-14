package logging

import "github.com/sirupsen/logrus"

func logFields(handler string) logrus.Fields {
	return logrus.Fields{
		"handler": handler,
	}
}

func LogError(handler string, err error) {
	logrus.WithFields(logFields(handler)).Error(err)
}
