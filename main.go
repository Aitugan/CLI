package main

import (
	proxymethod "CLI/proxyMethod"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	filename := "defaultConfig.json"
	app := cli.NewApp()
	app.Name = "CLI"
	app.Usage = "Idk, just a cli"
	app.Description = "For support, contact ..email.."
	app.UsageText = ""
	app.Version = "1.0.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Aitugan Mir",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "configFile",
			Value: "defaultConfig.json",
			Usage: "To pass configuration file (if none is provided, default file will be used)",
		},
	}
	app.Action = func(c *cli.Context) error {
		filename = c.String("configFile")
		return nil
	}

	app.Commands = []cli.Command{

		{
			Name:  "run",
			Usage: "Run the proxy server",
			Action: func(c *cli.Context) error {
				fmt.Println("Server is running\n")
				fmt.Printf("%v will be used as configuration file. If you changed your opinion, use reload command\n\n", filename)
				proxymethod.RunServer(filename)
				return nil
			},
		},
		{
			Name:  "reload",
			Usage: "Reload the proxy server. Stops previous one and opens another with new configuration file",
			Action: func(c *cli.Context) error {
				fmt.Println("Server is reloaded")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	Check(err)

}
