package utils

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

func Listen(uri *url.URL) (net.Listener, error) {
	var listener net.Listener
	switch uri.Scheme {
	case "unix":
		sock := getUnixPath(uri)
		os.Remove(sock)
		if l, err := net.Listen("unix", sock); err != nil {
			return nil, err
		} else {
			listener = l
		}
	case "http":
		if l, err := net.Listen("tcp", uri.Host); err != nil {
			return nil, err
		} else {
			listener = l
		}
	default:
		return nil, errors.New(fmt.Sprintf("unsupported target: %s", uri.String()))
	}
	return listener, nil
}

func getUnixPath(uri *url.URL) string {
	v, _ := strings.CutPrefix(uri.String(), "unix://")
	return v
}

func DialTimeout(uri *url.URL, timeout time.Duration) (net.Conn, error) {
	switch uri.Scheme {
	case "unix":
		path := getUnixPath(uri)
		return net.DialTimeout("unix", path, timeout)
	case "http":
		return net.DialTimeout("tcp", uri.Host, timeout)
	default:
		return nil, errors.New(fmt.Sprintf("unsupported target: %s", uri.String()))
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
