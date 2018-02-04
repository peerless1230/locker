package main

import (
	"fmt"

	"github.com/peerless1230/locker/container"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	//Usage: "Create a Docker-liked container\nlocker run -it [command],
	Usage: `Create a Docker-liked container
		locker run -[option, -option] [command]`,

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
		cli.BoolFlag{
			Name:  "m",
			Usage: "Memory limit",
		},
		cli.BoolFlag{
			Name:  "memory",
			Usage: "Memory limit",
		},
		cli.BoolFlag{
			Name:  "cpu-shares",
			Usage: "CPU shares (relative weight)",
		},
		cli.BoolFlag{
			Name:  "cpuset-cpus",
			Usage: " CPUs in which to allow execution (0-3, 0,1)",
		},
		cli.BoolFlag{
			Name:  "cpu-period",
			Usage: "Limit CPU CFS (Completely Fair Scheduler) period",
		},
		cli.BoolFlag{
			Name:  "cpu-quota",
			Usage: "Limit CPU CFS (Completely Fair Scheduler) quota",
		},
		cli.BoolFlag{
			Name:  "cpus",
			Usage: `Number of CPU. This is the equivalent of setting --cpu-period="100000" and --cpu-quota="n*100000"`,
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
