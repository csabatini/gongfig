package main

import (
	"github.com/urfave/cli"
	"os"
	"log"
	"fmt"
	"github.com/romanovskyj/gongfig/pkg/actions"
)

func main() {
	app := cli.NewApp()
	app.Name = "Gongfig"
	app.Usage = "Manage Kong configuration"
	app.Version = "0.0.1"

	flags := []cli.Flag {
		cli.StringFlag{
			Name: "url",
			Value: "http://localhost:8001",
			Usage: "Kong admin api url",
		},
		cli.StringFlag{
			Name: "file",
			Value: "config.yml",
			Usage: "File for export/import",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "export",
			Usage: "Obtain services and routes, write it to the config file",
			Action: func(c *cli.Context) error {
				fmt.Println("The configuration is exporting...")
				actions.Export(c.String("url"), c.String("file"))

				return nil
			},
			Flags: flags,
		},
		{
			Name: "import",
			Usage: "Apply services and routes from configuration file to the kong deployment",
			Action: func(c *cli.Context) error {
				fmt.Println("The configuration is importing...")

				return nil
			},
			Flags: flags,
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}