package app

import (
	"bakery/application/internal/adapters/httpserver"
	"bakery/application/internal/adapters/pgstorage"
	"bakery/application/internal/config"
	"bakery/application/internal/domain/bakery"
	"bakery/application/internal/ports"
	"bakery/application/internal/service"

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
	storage := pgstorage.New()
	mBakery := bakery.New(logger, storage)
	mHttpserver := httpserver.New(config.ReadConfig(), mBakery, logger)
	app.Httpserver = mHttpserver
	app.Storage = storage
	if inService {
		app.Service = service.New(mHttpserver, logger, storage, "bakery")
	}

}

func (app *App) Start() {
	if app.InService {
		app.Service.Start()
	} else {
		app.Storage.Start()
		app.Httpserver.Start()
	}

}

func (app *App) Stop() {

}
