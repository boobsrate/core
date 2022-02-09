package observer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	defaultOpenTimeout  = time.Second * 30
	defaultCloseTimeout = time.Second * 30
)

type ContextOpener interface {
	Open(ctx context.Context) error
}

type ContextCloser interface {
	Close(ctx context.Context) error
}

type Opener interface {
	Open() error
}

type OpenerFunc func() error

func (f OpenerFunc) Open() error {
	return f()
}

type CloserFunc func() error

func (f CloserFunc) Close() error {
	return f()
}

type ContextCloserFunc func(ctx context.Context) error

func (f ContextCloserFunc) Close(ctx context.Context) error {
	return f(ctx)
}

type Closer interface {
	Close() error
}

type Observer struct {
	openTimeout  time.Duration
	closeTimeout time.Duration

	contextOpeners []ContextOpener
	contextClosers []ContextCloser
	openers        []Opener
	closers        []Closer
	uppers         []func(ctx context.Context)

	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

func NewObserver() *Observer {
	ctx, cancel := context.WithCancel(context.Background())

	observer := &Observer{
		openTimeout:  defaultOpenTimeout,
		closeTimeout: defaultCloseTimeout,
		ctx:          ctx,
		cancelFunc:   cancel,
	}
	observer.AddUpper(handleSignals)

	return &Observer{
		openTimeout:  defaultOpenTimeout,
		closeTimeout: defaultCloseTimeout,
		ctx:          ctx,
		cancelFunc:   cancel,
	}
}

func (o *Observer) AddContextOpener(contextOpener ContextOpener) {
	o.contextOpeners = append(o.contextOpeners, contextOpener)
}

func (o *Observer) AddContextCloser(contextCloser ContextCloser) {
	o.contextClosers = append(o.contextClosers, contextCloser)
}

func (o *Observer) AddOpener(opener Opener) {
	o.openers = append(o.openers, opener)
}

func (o *Observer) AddCloser(closer Closer) {
	o.closers = append(o.closers, closer)
}

func (o *Observer) AddUpper(upper func(ctx context.Context)) {
	o.uppers = append(o.uppers, upper)
}

func (o *Observer) Run() error {
	for _, opener := range o.openers {
		err := opener.Open()
		if err != nil {
			return fmt.Errorf("open %#v: %v", opener, err)
		}
	}

	openCtx, cancel := context.WithTimeout(context.Background(), o.openTimeout)
	defer cancel()

	for _, opener := range o.contextOpeners {
		if err := opener.Open(openCtx); err != nil {
			return fmt.Errorf("open with context %#v: %v", opener, err)
		}
	}

	for _, upper := range o.uppers {
		o.wg.Add(1)
		upper := upper

		go func() {
			defer o.wg.Done()
			defer o.cancelFunc()
			upper(o.ctx)
		}()
	}
	o.wg.Wait()

	cancel()
	err := o.stop()
	return err
}

func handleSignals(ctx context.Context) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigs)

	select {
	case <-sigs:
	case <-ctx.Done():
	}
}

func (o *Observer) stop() error {
	closeCtx, cancel := context.WithTimeout(context.Background(), o.closeTimeout)
	defer cancel()

	for _, closer := range o.closers {
		if err := closer.Close(); err != nil {
			return fmt.Errorf("close %#v: %v", closer, err)
		}
	}

	for _, closer := range o.contextClosers {
		if err := closer.Close(closeCtx); err != nil {
			return fmt.Errorf("close with context %#v: %v", closer, err)
		}
	}
	return nil
}
