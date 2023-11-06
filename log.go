package suda

import (
	"io"

	"dxkite.cn/log"
)

type moduleLogger struct {
	name string
}

func (m moduleLogger) Write(b []byte) (int, error) {
	log.Debug("[" + m.name + "] " + string(b))
	return len(b), nil
}

func makeLoggerWriter(name string) io.Writer {
	return moduleLogger{name: name}
}
