package app

import (
	"context"
	"errors"
	"go-serv/constant"
	"go-serv/server"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type ShotdownCallback func(c context.Context) error

type Option func(app *App)

type App struct {
	servers []server.Service
	cbs     []ShotdownCallback
	ctx     context.Context
	cancel  func()
	opts    options
	once    sync.Once
	sig     chan os.Signal
	wg      sync.WaitGroup
}

func WithShotdownCallback(cbs ...ShotdownCallback) Option {
	return func(app *App) {
		app.cbs = cbs
	}
}

func Wait(ctx context.Context, done chan struct{}) error {
	select {
	case <-ctx.Done():
		return errors.New(constant.TimeOutErr)
	case <-done:
		return nil
	}

}

func StoreCacheToDBCallback(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		log.Println("将缓存数据重写回mysql")
		time.Sleep(3 * time.Second)
		log.Println("重写完成")
		done <- struct{}{}
	}()
	return Wait(ctx, done)
}

func ReleaseResourceCallback(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		log.Println("释放应用资源!")
		time.Sleep(3 * time.Second)
		log.Println("释放完成!")
		done <- struct{}{}
	}()
	return Wait(ctx, done)
}

func (a *App) InitApp() {
	a.once.Do(func() {
		a.servers = server.BuidServers()
		cbs := []ShotdownCallback{StoreCacheToDBCallback, ReleaseResourceCallback}
		WithShotdownCallback(cbs...)(a)

		opts := options{
			ctx:         context.Background(),
			stopTimeout: constant.TimeOut * time.Second,
		}

		switch runtime.GOOS {
		case "windows":
			opts.sigs = []os.Signal{os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGABRT, syscall.SIGTERM}
		case "linux":
			//  还有一个 syscall.SIGSTOP,
			opts.sigs = []os.Signal{os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGABRT, syscall.SIGFPE, syscall.SIGSEGV, syscall.SIGTERM}
		case "darwin":
			// 还有一个 syscall.SIGSTOP
			opts.sigs = []os.Signal{os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGABRT, syscall.SIGTERM}
		}

		ctx, cancel := context.WithCancel(opts.ctx)
		a.sig = make(chan os.Signal, 1)
		a.ctx = ctx
		a.opts = opts
		a.cancel = cancel
		a.wg = sync.WaitGroup{}
	})

}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	for _, showdown := range a.cbs {
		err := showdown(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Run() error {
	a.InitApp()
	eg, ctx := errgroup.WithContext(a.ctx)
	for _, srv := range a.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done()
			sctx, cacel := context.WithCancel(a.opts.ctx)
			defer cacel()
			return srv.Stop(sctx)
		})
		a.wg.Add(1)
		eg.Go(func() error {
			a.wg.Done()
			return srv.Serve(ctx)

		})
	}
	signal.Notify(a.sig, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-a.sig:
			err := a.Stop()
			a.wg.Add(1)
			go func() {
				select {
				case <-a.sig:
					os.Exit(1)
				case <-time.After(a.opts.stopTimeout):
					os.Exit(1)
				}
				a.wg.Done()
			}()
			return err
		}
	})
	a.wg.Wait()
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func NewApp() *App {
	return &App{}
}
