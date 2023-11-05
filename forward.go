package suda

import (
	"io"
)

func transport(src, dst io.ReadWriteCloser) (up, down int64, err error) {
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
