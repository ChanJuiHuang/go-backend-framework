package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

type SignalCallback struct {
	Signals      []os.Signal
	CallbackFunc func(ctx context.Context)
}

type App struct {
	wg                  sync.WaitGroup
	startingCallbacks   []func()
	startedCallbacks    []func()
	signalCallbacks     []SignalCallback
	asyncCallbacks      []func()
	terminatedCallbacks []func()
}

func New(
	startingCallbacks []func(),
	startedCallbacks []func(),
	signalCallbacks []SignalCallback,
	asyncCallbacks []func(),
	terminatedCallbacks []func(),
) *App {
	return &App{
		wg:                  sync.WaitGroup{},
		startingCallbacks:   startingCallbacks,
		startedCallbacks:    startedCallbacks,
		signalCallbacks:     signalCallbacks,
		asyncCallbacks:      asyncCallbacks,
		terminatedCallbacks: terminatedCallbacks,
	}
}

func (app *App) runStartingCallbacks() {
	for _, startingCallback := range app.startingCallbacks {
		startingCallback()
	}
}

func (app *App) runStartedCallbacks() {
	for _, startedCallback := range app.startedCallbacks {
		startedCallback()
	}
}

func (app *App) runSignalCallBacks() {
	app.wg.Add(len(app.signalCallbacks))

	for _, signalCallback := range app.signalCallbacks {
		ctx, stop := signal.NotifyContext(context.Background(), signalCallback.Signals...)
		go func() {
			<-ctx.Done()
			stop()
		}()

		callback := signalCallback.CallbackFunc
		go func() {
			defer app.wg.Done()
			callback(ctx)
		}()
	}
}

func (app *App) runAsyncCallbacks() {
	for _, asyncCallback := range app.asyncCallbacks {
		go asyncCallback()
	}
}

func (app *App) runTerminatedCallbacks() {
	for _, terminatedCallback := range app.terminatedCallbacks {
		terminatedCallback()
	}
}

func (app *App) Run(executerFunc func()) {
	app.runStartingCallbacks()
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		executerFunc()
	}()
	app.runStartedCallbacks()
	app.runSignalCallBacks()
	app.runAsyncCallbacks()

	app.wg.Wait()
	app.runTerminatedCallbacks()
}
