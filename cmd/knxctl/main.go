package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/dpt"
	"github.com/vapourismo/knx-go/knx/knxnet"
)

func discover() error {
	client, err := knx.DescribeTunnel(fmt.Sprintf("%s:%s", server, port), time.Millisecond*750)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", client)

	return nil
}

func listen() error {
	for {
		tunnel, err := knx.NewTunnel(fmt.Sprintf("%s:%s", server, port), knxnet.TunnelLayerData, knx.DefaultTunnelConfig)
		if err != nil {
			fmt.Printf("Error while creating: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		for {
			select {
			case msg, open := <-tunnel.Inbound():
				if !open {
					fmt.Println("tunnel channel closed")
					return errors.New("tunnel channel closed")
				}

				if ind, ok := msg.(*cemi.LDataInd); ok {
					fmt.Printf("%+v", ind)
				}
			}
		}
	}
}

func sendbool() error {
	client, err := knx.NewGroupTunnel(fmt.Sprintf("%s:%s", server, port), knx.DefaultTunnelConfig)
	if err != nil {
		fmt.Printf("Error while creating: %v\n", err)
		return err
	}
	defer client.Close()

	gd, err := cemi.NewGroupAddrString(group)
	if err != nil {
		return err
	}

	err = client.Send(knx.GroupEvent{
		Command:     knx.GroupWrite,
		Destination: gd,
		Data:        dpt.DPT_1001(value).Pack(),
	})
	if err != nil {
		fmt.Printf("Error while sending: %v\n", err)
		return err
	}

	for msg := range client.Inbound() {
		var temp dpt.DPT_1001

		err := temp.Unpack(msg.Data)
		if err != nil {
			continue
		}

		fmt.Printf("%+v: %v", msg, temp)
		break
	}
	return nil
}

var (
	server string
	port   string
	group  string
	value  bool
)

func main() {
	app := &cli.App{
		Name:  "knxctl",
		Usage: "knxctl [action]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "server",
				Aliases:     []string{"s"},
				Value:       "127.0.0.1",
				Usage:       "server IP Address",
				Destination: &server,
			},
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       "3671",
				Usage:       "server Port",
				Destination: &port,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "discover",
				Usage: "discover KNX server",
				Action: func(cCtx *cli.Context) error {
					return discover()
				},
			},
			{
				Name:  "listen",
				Usage: "listen KNX messages",
				Action: func(cCtx *cli.Context) error {
					return listen()
				},
			},
			{
				Name:  "send-bool",
				Usage: "send a boolean KNX message",
				Action: func(cCtx *cli.Context) error {
					return sendbool()
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "group",
						Aliases:     []string{"g"},
						Value:       "1",
						Usage:       "KNX group",
						Destination: &group,
					},
					&cli.BoolFlag{
						Name:        "value",
						Aliases:     []string{"v"},
						Value:       true,
						Usage:       "KNX value",
						Destination: &value,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
