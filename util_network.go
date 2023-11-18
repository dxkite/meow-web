package suda

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Port struct {
	Type string   `yaml:"type"`
	Unix UnixPort `yaml:"unix"`
	Http HttpPort `yaml:"http"`
}

type UnixPort struct {
	Path string `yaml:"path"`
}

type HttpPort struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (port Port) String() string {
	switch port.Type {
	case "unix":
		return fmt.Sprintf("unix://%s", port.Unix.Path)
	case "http":
		return fmt.Sprintf("http://%s:%d", port.Http.Host, port.Http.Port)
	default:
		return port.Type
	}
}

func Listen(port Port) (net.Listener, error) {
	var listener net.Listener
	switch port.Type {
	case "unix":
		sock := port.Unix.Path
		os.Remove(sock)
		if l, err := net.Listen("unix", sock); err != nil {
			return nil, err
		} else {
			listener = l
		}
	case "http":
		addr := fmt.Sprintf("%s:%d", port.Http.Host, port.Http.Port)
		if l, err := net.Listen("tcp", addr); err != nil {
			return nil, err
		} else {
			listener = l
		}
	default:
		return nil, errors.New(fmt.Sprintf("unsupported target: %s", port.String()))
	}
	return listener, nil
}

func DialTimeout(port Port, timeout time.Duration) (net.Conn, error) {
	switch port.Type {
	case "unix":
		return net.DialTimeout("unix", port.Unix.Path, timeout)
	case "http":
		addr := fmt.Sprintf("%s:%d", port.Http.Host, port.Http.Port)
		return net.DialTimeout("tcp", addr, timeout)
	default:
		return nil, errors.New(fmt.Sprintf("unsupported target: %s", port.String()))
	}
}

func Transport(src, dst io.ReadWriter) (up, down int64, err error) {
	var closeCh = make(chan struct{})
	var errCh = make(chan error)

	go func() {
		// remote -> local
		var _err error
		if down, _err = io.Copy(src, dst); _err != nil && _err != io.EOF {
			errCh <- _err
			return
		}
		closeCh <- struct{}{}
	}()

	go func() {
		// local -> remote
		var _err error
		if up, _err = io.Copy(dst, src); _err != nil && _err != io.EOF {
			errCh <- _err
			return
		}
		closeCh <- struct{}{}
	}()

	select {
	case err = <-errCh:
		return
	case <-closeCh:
		return
	}
}
