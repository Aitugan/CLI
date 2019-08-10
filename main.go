package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Interface string     `json:"interface"`
	Upstreams []Upstream `json:"upstreams"`
}

type Upstream struct {
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	Backends    []string `json:"backends"`
	ProxyMethod string   `json:"proxyMethod"`
}

func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	Check(err)
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Some data Some data Some data Some data Some data")
}

func main() {
	filename := "defaultConfig.json"
	config, err := LoadConfiguration(filename)
	Check(err)
	fmt.Println(config.Upstreams[0].Path)

	app := cli.NewApp()
	app.Name = "CLI"
	// app.Usage = "Idk, just a cli"
	// app.Description = "For support, contact ..email.."
	// app.UsageText = ""
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
		// if c.NArg() > 0 {
		filename = c.String("configFile") //c.Args().Get(0)
		// }
		// fmt.Println(c.Args().Get(0))
		// fmt.Println(c.String("configFile"))
		// fmt.Println(filename)
		// fmt.Println("The file ", filename, " will be used as a config file")
		return nil
	}
	// config := Config{}
	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "Run the proxy server",
			Action: func(c *cli.Context) error {
				fmt.Println("Server is running\n")
				fmt.Printf("%v will be used as configuration file. If you changed your opinion, use reload command\n\n", filename)
				//TODO
				//Start a new proxy server
				//open in browser... I suppose

				app.Flags[0]

				http.HandleFunc("/", handler)
				http.ListenAndServe(config.Interface, nil)

				return nil
			},
		},
		// {
		// 	Name: "reload",
		// 	Usage: "Reload the proxy server. Stops previous one and opens another with new configuration file"
		// 	Action: func(c *cli.Context) error {
		// 		fmt.Println("Server is reloaded")
		// 		return nil
		// 	}
		// }
	}

	err = app.Run(os.Args)
	Check(err)
}
