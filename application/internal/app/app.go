package app

import (
	"bakery/application/internal/adapters/httpserver"
	"bakery/application/internal/adapters/pgstorage"
	"bakery/application/internal/config"
	"bakery/application/internal/domain/bakery"
	"bakery/application/internal/ports"
	"bakery/application/internal/service"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	Service    *service.MyService
	Httpserver *httpserver.Server
	Storage    ports.Storage
	InService  bool
}

func New(inService bool, log *logrus.Logger) *App {
	app := App{
		InService: inService,
	}
	logger := logrus.New()
	config := config.ReadConfig()
	storage := pgstorage.New(logger, config)
	mBakery := bakery.New(logger, storage)
	mHttpserver := httpserver.New(config, mBakery, logger)
	app.Httpserver = mHttpserver
	app.Storage = storage
	if inService {
		app.Service = service.New(mHttpserver, logger, storage, "bakery")
	}
	return &app
}

func (app *App) Start() {
	if app.InService {
		app.Service.Start(context.Background())
	} else {
		ctx, cancelFunc := context.WithCancel(context.Background())
		storageStarted := make(chan int)
		go app.Storage.Start(ctx, storageStarted)
		<-storageStarted
		go app.Httpserver.Start(ctx)
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c
		cancelFunc()
		<-time.After(10 * time.Second)
	}

}
