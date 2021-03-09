package main

import (
	"bakery/pkg/service"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/sys/windows/svc"
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

	const svcName = "bakery"

	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}
	if inService {
		service.RunService(svcName)
		return
	}

	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "run":
		service.StartServer()
		for {
			//		time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
		}
	case "install":
		err = service.InstallService(svcName, "my service")
	case "remove":
		err = service.RemoveService(svcName)
	case "start":
		err = service.StartService(svcName)
	case "stop":
		err = service.ControlService(svcName, svc.Stop, svc.Stopped)
	case "pause":
		err = service.ControlService(svcName, svc.Pause, svc.Paused)
	case "continue":
		err = service.ControlService(svcName, svc.Continue, svc.Running)
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}
	return

}
