package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.Out = os.Stdout
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.InfoLevel)
}
