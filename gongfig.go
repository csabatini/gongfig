package main

import (
	"github.com/urfave/cli"
	"os"
	"log"
	"fmt"
	"github.com/csabatini/gongfig/pkg/actions"
)

func getApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Gongfig"
	app.Usage = "Manage Kong configuration"
	app.Version = "0.0.1"

	flags := []cli.Flag {
		cli.StringFlag{
			Name: "url",
			Value: actions.DefaultURL,
			Usage: "Kong admin api url",
		},
		cli.StringFlag{
			Name: "file",
			Value: "config.yml",
			Usage: "File for export/import",
		},
		cli.StringFlag{
			Name: "key",
			Value: "",
			Usage: "Optional key for admin API with auth plugin enabled",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "export",
			Usage: "Obtain services and routes, write it to the config file",
			Action: func(c *cli.Context) error {
				fmt.Println("The configuration is exporting...")
				actions.Export(c.String("url"), c.String("file"), c.String("key"))

				return nil
			},
			Flags: flags,
		},
		{
			Name: "import",
			Usage: "Apply services and routes from configuration file to the kong deployment",
			Action: func(c *cli.Context) error {
				fmt.Println("The configuration is importing...")
				actions.Import(c.String("url"), c.String("file"))

				return nil
			},
			Flags: flags,
		},
		{
			Name: "flush",
			Usage: "Delete all services and routes from configuration file to the kong deployment",
			Action: func(c *cli.Context) error {
				actions.Flush(c.String("url"), c.String("key"))

				return nil
			},
			Flags: flags,
		},
	}

	return app
}

func main() {
	app := getApp()

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
