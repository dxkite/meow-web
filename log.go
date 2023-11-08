package suda

import (
	"bufio"
	"bytes"
	"io"

	"dxkite.cn/log"
)

type moduleLogger struct {
	name string
}

func (m moduleLogger) Write(b []byte) (int, error) {
	s := bufio.NewScanner(bytes.NewReader(b))
	for s.Scan() {
		t := s.Text()
		log.Debug("[" + m.name + "] " + t)
	}
	return len(b), nil
}

func MakeNameLoggerWriter(name string) io.Writer {
	return moduleLogger{name: name}
}
