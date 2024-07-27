package log

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var G = GetLogger

var Default = &Entry{
	Logger: logrus.StandardLogger(),
}

type Fields = map[string]interface{}

type Entry = logrus.Entry

type Level = logrus.Level

const (
	TraceLevel Level = logrus.TraceLevel
	DebugLevel Level = logrus.DebugLevel
	InfoLevel  Level = logrus.InfoLevel
	WarnLevel  Level = logrus.WarnLevel
	ErrorLevel Level = logrus.ErrorLevel
	FatalLevel Level = logrus.FatalLevel
	PanicLevel Level = logrus.PanicLevel
)

func SetLevel(level string) error {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	Default.Logger.SetLevel(lv)
	return nil
}

func GetLevel() Level {
	return Default.Logger.GetLevel()
}

type OutputFormat string

const (
	TextFormat OutputFormat = "text"
	JSONFormat OutputFormat = "json"
)

func SetFormat(format OutputFormat) error {
	switch format {
	case TextFormat:
		Default.Logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FullTimestamp:   true,
		})
		return nil
	case JSONFormat:
		Default.Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
		return nil
	default:
		return fmt.Errorf("unknown log format: %s", format)
	}
}

func GetLogger(ctx context.Context) *Entry {
	return Default.WithContext(ctx)
}
