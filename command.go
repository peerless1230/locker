package main

import (
	"fmt"

	"./container"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	//Usage: "Create a Docker-liked container\nlocker run -it [command],
	Usage: `Create a Docker-liked container
		locker run -it [command]`,

	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "t",
			Usage: "enable tty for container",
		},
		cli.BoolFlag{
			Name:  "i",
			Usage: "enable interactive to container",
		},
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable interactive tty to container",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		tty := (context.Bool("i") && context.Bool("t") || context.Bool("it"))

		if context.Args().Get(1) != "" {
			Run(tty, cmd, context.Args()[1:])
		}
		Run(tty, cmd, []string{})
		return nil

	},
}

var initCommand = cli.Command{
	Name: "init",
	Usage: `Init the Docker-liked container and run user's commands
			locker init [command]`,
	Action: func(context *cli.Context) error {
		log.Infof("init started:")
		args := context.Args().Get(0)
		log.Infof("Command: %s", args)
		err := container.NewContainerInitProcess(args)
		return err
	},
}
