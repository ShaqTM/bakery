package service

import (
	"bakery/application/api"
	"bakery/application/store"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

type MyService struct{ Log *logrus.Logger }

func (m *myservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	go StartServer(m.Log)

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				// golang.org/x/sys/windows/svc.TestExample is verifying this output.
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

// RunService - Запуск сервиса
func RunService(name string, log *logrus.Logger) {
	var err error
	elog, err = eventlog.Open(name)
	if err != nil {
		return
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	log.Info(fmt.Sprintf("starting %s service", name))

	run := svc.Run
	err = run(name, &myservice{Log: log})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
	log.Info(fmt.Sprintf("%s service stopped", name))
}

// StartServer - запуск web-сервера
func StartServer(log *logrus.Logger) {
	var db *sql.DB
	mdb := store.MDB{Pdb: &db, Log: log}
	(&mdb).InitDatabase()
	router := mux.NewRouter()
	api.AddMaterialsRoutes(&router, mdb)
	api.AddUnitsRoutes(&router, mdb)
	api.AddOrdersRoutes(&router, mdb)
	api.AddRecipesRoutes(&router, mdb)

	http.ListenAndServe(":5000", router)

}
