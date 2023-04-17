package main

import (
	"flag"
	"fmt"
	"goCrunchQueue/internal"
	"os"
	"os/signal"
)

var (
	version   string // set version
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

func main() {
	configPath := flag.String("config", "./config/config.yaml", "path to the config file")
	flag.Parse()

	var formattedVersion = fmt.Sprintf("Goshboard Version %v (Git %v | build at %v)", version, sha1ver, buildTime)

	internal.HandleFatalErrorf(internal.ValidateSetup(), "Application Pre-Check Failed")
	addExitHandler()
	internal.StartApplication(formattedVersion, configPath)
}

func addExitHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			internal.StopApplication()
			os.Exit(1)
		}
	}()
}
