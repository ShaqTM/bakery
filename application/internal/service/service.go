package service

import (
	"bakery/application/internal/adapters/httpserver"
	"bakery/application/internal/ports"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

type MyService struct {
	HttpServer *httpserver.Server
	Log        *logrus.Logger
	Storage    ports.Storage
	Name       string
}

func New(httpServer *httpserver.Server, log *logrus.Logger, storage ports.Storage, name string) *MyService {
	myService := MyService{
		HttpServer: httpServer,
		Log:        log,
		Storage:    storage,
		Name:       name,
	}
	return &myService
}

func (m *MyService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	go m.Start()

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
func (m *MyService) RunService() {
	var err error
	elog, err = eventlog.Open(m.Name)
	if err != nil {
		return
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", m.Name))
	m.Log.Info(fmt.Sprintf("starting %s service", m.Name))

	run := svc.Run
	err = run(m.Name, m)
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", m.Name))
	m.Log.Info(fmt.Sprintf("%s service stopped", m.Name))
}

// StartServer - запуск web-сервера
func (m *MyService) Start() {
	m.Storage.Start()
	m.HttpServer.Start()

}
