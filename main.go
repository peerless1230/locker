package main

import (
	"os"
	"sort"

	"github.com/peerless1230/locker/common"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "locker"
	app.Version = "v1.0.0-alpha"
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
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	for _, ele := range app.Commands {
		sort.Sort(cli.FlagsByName(ele.Flags))
	}
	common.CheckError(app.Run(os.Args))
}
