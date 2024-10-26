package closing

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Func func(context.Context) error

type Closer struct {
	mux sync.Mutex
	fns []Func
}

func New() *Closer {
	return &Closer{}
}

func (c *Closer) Add(fn Func) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.fns = append(c.fns, fn)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	var (
		wg       sync.WaitGroup
		complete = make(chan struct{}, 1)

		mux  sync.Mutex
		errs error
	)

	go func() {
		wg.Wait()
		complete <- struct{}{}
	}()

	for _, fn := range c.fns {
		wg.Add(1)
		go func(fn Func) {
			defer wg.Done()
			err := fn(ctx)

			mux.Lock()
			errs = errors.Join(errs, err)
			mux.Unlock()
		}(fn)
	}

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return errors.New("closing: context deadline exceeded or canceled")
	}

	if errs != nil {
		return fmt.Errorf("closing: %w", errs)
	}

	return nil
}

func WaitQuit() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}
