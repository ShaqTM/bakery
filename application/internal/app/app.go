package app

import (
	"bakery/application/internal/adapters/httpserver"
	"bakery/application/internal/adapters/pgstorage"
	"bakery/application/internal/config"
	"bakery/application/internal/domain/bakery"
	"bakery/application/internal/ports"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	Httpserver *httpserver.Server
	Storage    ports.Storage
}

func New(log *logrus.Logger) *App {
	app := App{}
	logger := logrus.New()
	config := config.ReadConfig()
	storage := pgstorage.New(logger, config)
	mBakery := bakery.New(logger, storage)
	mHttpserver := httpserver.New(config, mBakery, logger)
	app.Httpserver = mHttpserver
	app.Storage = storage
	return &app
}

func (app *App) Start() {
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
