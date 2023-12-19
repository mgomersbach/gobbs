package logger

import (
	"gobbs/config"
	"os"

	"github.com/sirupsen/logrus"
)

// InitializeLogger initializes the logger based on the provided configuration
func InitializeLogger(cfg *config.Config) *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	// Set up additional logger configurations here
	return log
}
