package log

import (
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	})

	// Log output:
	// {"file":"example.go","func":"example","level":"info","msg":"","time":"2006-01-02T15:04:05.000Z"}
}
