package logger

import (
	"context"
	"fmt"
	"io"
	"oauth-server/utils"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

var logger *logrus.Logger

type LogPrintln struct {
	Ctx       context.Context
	FileName  string
	FuncName  string
	TraceData interface{}
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

	//
	log.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	)))

	logger = log
}

func GetLogger() *logrus.Logger {
	return logger
}

func Println(params LogPrintln) {
	GetLogger().WithContext(params.Ctx).WithFields(logrus.Fields{
		"filename":   params.FileName,
		"func_name":  params.FuncName,
		"trace_data": params.TraceData,
	}).Println(params.Msg)
}
