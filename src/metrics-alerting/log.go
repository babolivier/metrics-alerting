package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

type utcFormatter struct {
	logrus.Formatter
}

func (f utcFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Time = entry.Time.UTC()
	return f.Formatter.Format(entry)
}

func logConfig() error {
	var disableColors = false

	if len(*logFile) > 0 {
		f, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE, 0655)
		if err != nil {
			return err
		}
		logrus.SetOutput(f)

		disableColors = true
	}

	logrus.SetFormatter(&utcFormatter{
		&logrus.TextFormatter{
			TimestampFormat:  "2006-01-02T15:04:05.000000000Z07:00",
			FullTimestamp:    true,
			DisableColors:    disableColors,
			DisableTimestamp: false,
			DisableSorting:   false,
		},
	})

	return nil
}
