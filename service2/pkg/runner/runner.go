package runner

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type ServiceRunner interface {
	Start(ctx context.Context) error
	ShutDown(ctx context.Context) error
}

type Runner struct {
	mainService ServiceRunner
	services    []ServiceRunner
}

func New(main ServiceRunner, services ...ServiceRunner) *Runner {
	return &Runner{mainService: main, services: services}
}

func (r *Runner) start(ctx context.Context) {
	for _, service := range r.services {
		if err := service.Start(ctx); err != nil {
			log.Println("Services start error: ", err)
			return
		}
	}
	if err := r.mainService.Start(ctx); err != nil {
		log.Println("MainService start error ", err)
		return
	}
}

func (r *Runner) stop(ctx context.Context) {
	if err := r.mainService.ShutDown(ctx); err != nil {
		log.Println("MainService stop error ", err)
		return
	}
	for i := len(r.services) - 1; i >= 0; i-- {
		if err := r.services[i].ShutDown(ctx); err != nil {
			log.Println("Services stop error ", err)
			return
		}
	}
}

func (r *Runner) RunnerRun() {
	ctx := context.Background()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	r.start(ctx)
	<-ch
	r.stop(ctx)
}
