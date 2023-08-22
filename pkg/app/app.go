package app

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
)

type StartingCallback func(ctx context.Context, wg *sync.WaitGroup)
type StartedCallback func()
type RunnerFunc func(wg *sync.WaitGroup)
type TerminatedCallback func()

type App struct {
	wg                  *sync.WaitGroup
	startingCallbacks   []StartingCallback
	startedCallbacks    []StartedCallback
	terminatedCallbacks []TerminatedCallback
	runnerFunc          RunnerFunc
}

func New(
	startingCallbacks []StartingCallback,
	startedCallbacks []StartedCallback,
	terminatedCallbacks []TerminatedCallback,
	runnerFunc RunnerFunc,
) *App {
	return &App{
		wg:                  &sync.WaitGroup{},
		startingCallbacks:   startingCallbacks,
		startedCallbacks:    startedCallbacks,
		terminatedCallbacks: terminatedCallbacks,
		runnerFunc:          runnerFunc,
	}
}

func (app *App) runStartingCallbacks() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	go func() {
		<-ctx.Done()
		stop()
	}()

	app.wg.Add(len(app.startingCallbacks))
	for _, callback := range app.startingCallbacks {
		go callback(ctx, app.wg)
	}
}

func (app *App) runStartedCallback() {
	for _, startedCallback := range app.startedCallbacks {
		go startedCallback()
	}
}

func (app *App) runTerminatedCallback() {
	for _, terminatedCallback := range app.terminatedCallbacks {
		terminatedCallback()
	}
}

func (app *App) Run() {
	app.runStartingCallbacks()
	app.wg.Add(1)
	go app.runnerFunc(app.wg)
	app.runStartedCallback()

	app.wg.Wait()
	app.runTerminatedCallback()
}
