// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix pilot service
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/urfave/cli"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/types"
	"openpitrix.io/openpitrix/pkg/service/pilot"
	"openpitrix.io/openpitrix/pkg/service/pilot/pilotutil"
)

func main() {
	app := cli.NewApp()
	app.Name = "openpitrix-pilot"
	app.Usage = "openpitrix-pilot provides pilot service."
	app.Version = "0.0.0"

	app.UsageText = `openpitrix-pilot [global options] command [options] [args...]

EXAMPLE:
   openpitrix-pilot gen-config
   openpitrix-pilot info
   openpitrix-pilot list
   openpitrix-pilot ping
   openpitrix-pilot getv key
   openpitrix-pilot confd-start
   openpitrix-pilot serve
   openpitrix-pilot tour`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "pilot-config.json",
			Usage:  "pilot config file",
			EnvVar: "OPENPITRIX_PILOT_CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "gen-config",
			Usage: "gen default config",

			Action: func(c *cli.Context) {
				fmt.Println(pilot.NewDefaultConfigString())
				return
			},
		},

		{
			Name:  "info",
			Usage: "show pilot service info",

			Action: func(c *cli.Context) {
				cfg := pilotutil.MustLoadPilotConfig(c.GlobalString("config"))

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}
				defer conn.Close()

				info, err := client.GetPilotConfig(context.Background(), &pbtypes.Empty{})
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}

				fmt.Println(JSONString(info))
				return
			},
		},
		{
			Name:      "list",
			Usage:     "list frontgate nodes",
			ArgsUsage: "[regexp]",

			Action: func(c *cli.Context) {
				cfg := pilotutil.MustLoadPilotConfig(c.GlobalString("config"))

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}
				defer conn.Close()

				list, err := client.GetFrontdateList(context.Background(), &pbtypes.Empty{})
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}

				re := c.Args().First()
				for _, v := range list.GetIdList() {
					if re == "" {
						fmt.Println(v)
						continue
					}
					matched, err := regexp.MatchString(re, v)
					if err != nil {
						logger.Fatal(err)
						os.Exit(1)
					}
					if matched {
						fmt.Println(v)
					}
				}
				return
			},
		},

		{
			Name:  "ping",
			Usage: "ping pilot/frontgate/drone service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "endpoint-type",
					Value: "pilot",
					Usage: "set endpoint type (pilot/frontgate/drone)",
				},

				cli.StringFlag{
					Name:  "frontgate-id",
					Value: "openpitrix-frontgate-001",
				},
				cli.StringFlag{
					Name:  "drone-host",
					Value: "localhost",
				},
				cli.IntFlag{
					Name:  "drone-port",
					Value: constants.DroneServicePort,
				},
			},

			Action: func(c *cli.Context) {
				cfg := pilotutil.MustLoadPilotConfig(c.GlobalString("config"))

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}
				defer conn.Close()

				switch s := c.String("endpoint-type"); s {
				case "pilot":
					_, err = client.PingPilot(context.Background(), &pbtypes.Empty{})
					if err != nil {
						logger.Fatal(err)
						os.Exit(1)
					}
					fmt.Println("OK")
					return

				case "frontgate":
					_, err = client.PingFrontgate(context.Background(), &pbtypes.FrontgateId{
						Id: c.String("frontgate-id"),
					})
					if err != nil {
						logger.Fatal(err)
						os.Exit(1)
					}
					fmt.Println("OK")
					return

				case "drone":
					_, err = client.PingDrone(
						context.Background(),
						&pbtypes.DroneEndpoint{
							FrontgateId: c.String("frontgate-id"),
							DroneIp:     c.String("drone-host"),
							DronePort:   int32(c.Int("drone-port")),
						},
					)
					if err != nil {
						logger.Fatal(err)
						os.Exit(1)
					}
					fmt.Println("OK")
					return

				default:
					logger.Fatalf("unknown endpoint type: %s\n", s)
					os.Exit(1)
				}

				return
			},
		},

		{
			Name:  "confd-status",
			Usage: "get confd service status",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "frontgate-id",
					Value: "openpitrix-frontgate-001",
				},
				cli.StringFlag{
					Name:  "drone-host",
					Value: "",
				},
				cli.IntFlag{
					Name:  "drone-port",
					Value: constants.DroneServicePort,
				},
			},

			Action: func(c *cli.Context) {
				cfg := pilotutil.MustLoadPilotConfig(c.GlobalString("config"))

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}
				defer conn.Close()

				reply, err := client.IsConfdRunning(context.Background(), &pbtypes.DroneEndpoint{
					FrontgateId: c.String("frontgate-id"),
					DroneIp:     c.String("drone-host"),
					DronePort:   int32(c.Int("drone-port")),
				})
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}

				if reply.GetValue() {
					fmt.Printf("confd on frontgate(%s)/drone(%s:%d) is running\n",
						c.String("frontgate-id"), c.String("drone-host"), c.Int("drone-port"),
					)
				} else {
					fmt.Printf("confd on frontgate(%s)/drone(%s:%d) not running\n",
						c.String("frontgate-id"), c.String("drone-host"), c.Int("drone-port"),
					)
				}

				return
			},
		},
		{
			Name:  "confd-start",
			Usage: "start confd service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "frontgate-id",
					Value: "openpitrix-frontgate-001",
				},
				cli.StringFlag{
					Name:  "drone-host",
					Value: "",
				},
				cli.IntFlag{
					Name:  "drone-port",
					Value: constants.DroneServicePort,
				},
			},

			Action: func(c *cli.Context) {
				cfg := pilotutil.MustLoadPilotConfig(c.GlobalString("config"))

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}
				defer conn.Close()

				_, err = client.StartConfd(context.Background(), &pbtypes.DroneEndpoint{
					FrontgateId: c.String("frontgate-id"),
					DroneIp:     c.String("drone-host"),
					DronePort:   int32(c.Int("drone-port")),
				})
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}

				fmt.Println("Done")
				return
			},
		},
		{
			Name:  "confd-stop",
			Usage: "stop confd service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "frontgate-id",
					Value: "openpitrix-frontgate-001",
				},
				cli.StringFlag{
					Name:  "drone-host",
					Value: "",
				},
				cli.IntFlag{
					Name:  "drone-port",
					Value: constants.DroneServicePort,
				},
			},

			Action: func(c *cli.Context) {
				cfg := pilotutil.MustLoadPilotConfig(c.GlobalString("config"))

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}
				defer conn.Close()

				_, err = client.StopConfd(context.Background(), &pbtypes.DroneEndpoint{
					FrontgateId: c.String("frontgate-id"),
					DroneIp:     c.String("drone-host"),
					DronePort:   int32(c.Int("drone-port")),
				})
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}

				fmt.Println("Done")
				return
			},
		},

		{
			Name:  "get-cmd-status",
			Usage: "get cmd status",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "task-id",
					Value: "",
				},
			},

			Action: func(c *cli.Context) {
				cfg := pilotutil.MustLoadPilotConfig(c.GlobalString("config"))

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}
				defer conn.Close()

				status, err := client.GetSubtaskStatus(context.Background(), &pbtypes.SubTaskId{
					TaskId: c.String("task-id"),
				})
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}

				fmt.Println(status.Status)
				return
			},
		},

		{
			Name:  "serve",
			Usage: "run as pilot service",
			Action: func(c *cli.Context) {
				cfg := pilotutil.MustLoadPilotConfig(c.GlobalString("config"))
				pilot.Serve(cfg)
				return
			},
		},
		{
			Name:  "tour",
			Usage: "show more examples",
			Action: func(c *cli.Context) {
				fmt.Println(tourTopic)
			},
		},
	}

	app.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Fprintf(ctx.App.Writer, "not found '%v'!\n", command)
	}

	app.Run(os.Args)
}

func JSONString(m interface{}) string {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return ""
	}
	data = bytes.Replace(data, []byte("\n"), []byte("\r\n"), -1)
	return string(data)
}

func Atoi(s string, defaultValue int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultValue
}

const tourTopic = `
openpitrix-pilot gen-config

openpitrix-pilot info
openpitrix-pilot list

openpitrix-pilot getv /
openpitrix-pilot getv /key
openpitrix-pilot getv / /key

openpitrix-pilot confd-start
openpitrix-pilot confd-stop
openpitrix-pilot confd-status

openpitrix-pilot serve

GOOS=windows openpitrix-pilot list
LIBCONFD_GOOS=windows openpitrix-pilot list
`
