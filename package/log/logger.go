package logger

import (
	"fmt"
	"gin-boilerplate/utils"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func Init() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: utils.DATE_TIME_FORMAT,
		PrettyPrint:     true,
	})
	log.SetLevel(logrus.InfoLevel)

	// Set up a multi writer hook
	currentTime := time.Now()
	rootDir := utils.RootDir()
	fileName := fmt.Sprintf("%s/logs/log_%s.log", rootDir, currentTime.Format(utils.DATE_FORMAT))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error_open_log_file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

	logger = log
}

func GetLogger() *logrus.Logger {
	return logger
}
