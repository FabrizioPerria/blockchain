package logging

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

func LoggerFactory(filePath string) *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
		FullTimestamp:   true,
		DisableColors:   false,
	})

	parent := path.Dir(filePath)
	os.MkdirAll(parent, 0o755)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o755)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return nil
	}
	if os.Getenv("SKIP_STDOUT_LOG") == "true" {
		l.SetOutput(f)
	} else {
		l.SetOutput(io.MultiWriter(os.Stdout, f))
	}

	return l
}
