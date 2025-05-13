package logger

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

type Fields = logrus.Fields

func SetupLogger() {
	logging := logrus.New()
	logging.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logging.Info("logged initiated using logrus")
	Logger = logging
}
