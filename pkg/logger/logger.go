package logger

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

const LogLevel = logrus.DebugLevel

const StandardPermission = 0o777

// GetLogger Init initialize logger.
func GetLogger(logPath string, level int, printLogsToStdOut bool) (Logger, error) {
	log := logrus.New()

	// l.SetReportCaller(true)
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	log.SetLevel(LogLevel)

	if err := os.MkdirAll(path.Dir(logPath), StandardPermission); err != nil {
		return Logger{}, fmt.Errorf("failed to create log dir: %w", err)
	}

	logFile, err := os.Create(logPath)
	if err != nil {
		return Logger{}, fmt.Errorf("wailed to init logger error: %w", err)
	}

	mw := io.MultiWriter(logFile)
	if printLogsToStdOut {
		mw = io.MultiWriter(os.Stdout, logFile)
	}

	log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: false,
	})
	log.SetOutput(mw)
	log.SetLevel(logrus.Level(level))

	return Logger{log}, nil
}
