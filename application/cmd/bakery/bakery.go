package main

import (
	"bakery/application/internal/app"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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
	app := app.New(log)
	log.SetOutput(os.Stdout)
	log.Info("Started")
	app.Start()
}
