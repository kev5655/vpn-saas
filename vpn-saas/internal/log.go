package internal

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var verbose bool

func SetVerbose(v bool) {
	verbose = v
}

func GetVerbose() bool {
	return verbose
}

func LogVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Printf(format, args...)
	}
}

func InitLogger(verbose bool) {
	if verbose {
		logrus.SetOutput(os.Stdout)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Info("Logger initialized in verbose mode: outputting logs to console")
	} else {
		logrus.SetOutput(os.Stdout)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		logrus.SetLevel(logrus.InfoLevel)
		logrus.Info("Logger initialized: outputting logs to console")
	}
}
