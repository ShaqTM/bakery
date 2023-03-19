package app

import (
	"bakery/application/internal/ports"
	"bakery/application/internal/service"

	"github.com/gorilla/mux"
)

type App struct {
	Storage *ports.Storage
	Service *service.MyService
	Router  **mux.Router
}

func CreateApp() *App {

}

func (app *App) Start() {

}

func (app *App) Stop() {

}
