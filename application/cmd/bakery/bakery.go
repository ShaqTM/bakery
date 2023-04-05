package main

import (
	"bakery/application/internal/app"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       run, install, remove, debug, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	os.Exit(2)
}
func main() {

	log := logrus.New()
	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}
	app := app.New(inService, log)
	if inService {

		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		logfile := app.Service.Name + ".log"
		if err == nil {
			logfile = dir + `\` + logfile
		}
		fmt.Println(dir)
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			elog, elogErr := eventlog.Open(app.Service.Name)
			if elogErr != nil {
				return
			}
			elog.Error(1, fmt.Sprintf("error opening file: %v", err))
			elog.Close()
			return
		}
		defer f.Close()
		// Output to stderr instead of stdout, could also be a file.
		log.SetOutput(f)
		app.Start()
		return
	} else {
		log.SetOutput(os.Stdout)
	}
	log.Info("Started")
	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "run":
		app.Start()
		for {
			//		time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
		}
	case "install":
		err = app.Service.InstallService("Bakery backend service")
	case "remove":
		err = app.Service.RemoveService()
	case "start":
		err = app.Service.StartService()
		if err == nil {
			log.Info("Service stared")
		}
	case "stop":
		err = app.Service.ControlService(svc.Stop, svc.Stopped)
		if err == nil {
			log.Info("Service stopped")
		}
	case "pause":
		err = app.Service.ControlService(svc.Pause, svc.Paused)
		if err == nil {
			log.Info("Service paused")
		}
	case "continue":
		err = app.Service.ControlService(svc.Continue, svc.Running)
		if err == nil {
			log.Info("Service continued")
		}
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, app.Service.Name, err)
	}
	return

}
