package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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
	filename := "config.json"
	config, err := LoadConfiguration(filename)
	Check(err)
	fmt.Println(config.Upstreams[0].Path)

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
				// app.Flags[0]

				config, _ := LoadConfiguration("config.json")
				for i := 0; i < len(config.Upstreams); i++ {
					RoundRobinRunner(config.Interface, config.Upstreams[i].Path, config.Upstreams[i].Method, config.Upstreams[i])
				}

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

	err = app.Run(os.Args)
	Check(err)

}

///////////////////////////////////////////////////////////////
type UpstreamNumber struct {
	Upstream
	Count int
}

func RoundRobinRunner(interf, path, method string, upstream Upstream) {
	r := mux.NewRouter()
	upNum := UpstreamNumber{upstream, 0}
	r.HandleFunc(path, upNum.roundRobinHandle).Methods(method)
	http.ListenAndServe(interf, r)
}

func (upstream *UpstreamNumber) roundRobinHandle(w http.ResponseWriter, r *http.Request) {

	w.Write(sendRequest(upstream.Backends[upstream.Count]))

	upstream.Count = (upstream.Count + 1) % len(upstream.Backends)

}
func sendRequest(url string) []byte {

	response, err := http.Get(url)

	// Check(err)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	// Check(err)
	if err != nil {
		panic(err)
	}
	return body
}
