package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/yametech/ym-go-scaffold/pkg/ymctl/api/gogen"
	"github.com/yametech/ym-go-scaffold/pkg/ymctl/api/new"
	rpc "github.com/yametech/ym-go-scaffold/pkg/ymctl/rpc/cli"
	"os"
	"runtime"
)

var (
	buildVersion = "0.1.0"
	commands     = []cli.Command{
		{
			Name:  "api",
			Usage: "generate api related files",
			Subcommands: []cli.Command{
				{
					Name:  "new",
					Usage: "generate files for go micro service by default template",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "o",
							Usage: "the output dir",
						},
					},
					Action: new.CreateServiceCommand,
				},
				{
					Name:  "go",
					Usage: "generate files for go micro service",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
					},
					Action: gogen.GoCommand,
				},
			},
		},
		{
			Name:  "rpc",
			Usage: "generate rpc code",
			Subcommands: []cli.Command{
				{
					Name:   "new",
					Usage:  `generate rpc demo service`,
					Action: rpc.RPCNew,
				},
				{
					Name:  "proto",
					Usage: `generate rpc from proto`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "src, s",
							Usage: "the file path of the proto source file",
						},
						cli.StringSliceFlag{
							Name:  "proto_path, I",
							Usage: `native command of protoc, specify the directory in which to search for imports. [optional]`,
						},
						cli.StringFlag{
							Name:  "dir, d",
							Usage: `the target path of the code`,
						},
					},
					Action: rpc.RPC,
				},
			},
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "a cli tool to generate code"
	app.Version = fmt.Sprintf("%s %s/%s", buildVersion, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	// cli already print error messages
	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:", err)
	}
}
