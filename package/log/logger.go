package logger

import (
	"fmt"
	"io"
	"oauth-server/utils"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type LogPrintln struct {
	FileName  string
	FuncName  string
	TraceData string
	Msg       string
}

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

func Println(params LogPrintln) {
	paramArr := make([]string, 0)
	v := reflect.ValueOf(params)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != "" {
			paramArr = append(paramArr, fmt.Sprintf("%v", v.Field(i).Interface()))
		}
	}
	GetLogger().Infof("%s \n", strings.Join(paramArr, " - "))
}
