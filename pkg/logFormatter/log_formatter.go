package logFormatter

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 36
	gray    = 37
)

type CustomLogFormatter struct{}

func (f *CustomLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.printColored(b, entry, "02-01-2006 15:04:05.000")
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *CustomLogFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, timestampFormat string) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = green
	}
	levelText := strings.ToUpper(entry.Level.String())

	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	fmt.Fprintf(b, "%s \x1b[%dm%s\x1b[0m %s", entry.Time.Format(timestampFormat), levelColor, levelText, entry.Message)
}
