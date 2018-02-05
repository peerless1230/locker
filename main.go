package main

import (
	"os"

	"github.com/peerless1230/locker/common"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "locker"
	// app.EnableBashCompletion = true
	app.Usage = "Imitate Docker to trace it's mechanism"
	app.Commands = []cli.Command{
		runCommand,
		initCommand,
	}
	app.Before = func(context *cli.Context) error {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}
	common.CheckError(app.Run(os.Args))
}
