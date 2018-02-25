package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/peerless1230/locker/cgroups/subsystems"
	"github.com/peerless1230/locker/container"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	//Usage: "Create a Docker-liked container\nlocker run -it [command],
	Usage: `Create a Docker-liked container
		locker run -[option, -option] [command]`,
	UseShortOptionHandling: true,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "t",
			Usage: "enable tty for container",
		},
		cli.BoolFlag{
			Name:  "i",
			Usage: "enable interactive to container",
		},
		cli.StringSliceFlag{
			Name:  "volume, v",
			Usage: "Mount volume '/host/path:/container/path'",
		},
		cli.StringFlag{
			Name:  "memory, m",
			Usage: "Memory limit",
		},
		cli.StringFlag{
			Name:  "cpu-shares",
			Usage: "CPU shares (relative weight)",
		},
		cli.StringFlag{
			Name:  "cpuset-cpus",
			Usage: " CPUs in which to allow execution (0-3, 0,1)",
		},
		cli.StringFlag{
			Name:  "cpu-period",
			Usage: "Limit CPU CFS (Completely Fair Scheduler) period",
		},
		cli.StringFlag{
			Name:  "cpu-quota",
			Usage: "Limit CPU CFS (Completely Fair Scheduler) quota",
		},
		cli.StringFlag{
			Name:  "cpus",
			Usage: `Number of CPU. This is the equivalent of setting --cpu-period="100000" and --cpu-quota="n*100000"`,
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		tty := context.Bool("i") && context.Bool("t")

		var cmdArray []string
		// copy Context args
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		log.Debugf("Volumes: %v", context.StringSlice("v"))
		resLimits := subsystems.ResourceLimitConfig{
			MemoryLimits: context.String("m"),
			CPUS:         context.String("cpus"),
			CPUPeriod:    context.String("cpu-period"),
			CPUQuota:     context.String("cpu-quota"),
			CPUSet:       context.String("cpuset-cpus"),
			CPUShare:     context.String("cpu-shares"),
		}
		Run(tty, cmdArray, &resLimits)
		return nil

	},
}

var initCommand = cli.Command{
	Name: "init",
	Usage: `Init the Docker-liked container and run user's commands
			locker init [command]`,
	Action: func(context *cli.Context) error {
		log.Infof("Init process started:")
		err := container.NewContainerInitProcess()
		return err
	},
}
