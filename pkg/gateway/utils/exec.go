package utils

import (
	"sync"
)

type ExecChain []func() error

func (ch ExecChain) Run() error {
	errChan := make(chan error)
	doneChan := make(chan struct{})
	go func() {
		wg := &sync.WaitGroup{}
		for _, fn := range ch {
			wg.Add(1)
			go func(fn func() error) {
				defer wg.Done()
				if err := fn(); err != nil {
					errChan <- err
					return
				}
			}(fn)
		}
		wg.Wait()
		doneChan <- struct{}{}
	}()
	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}
