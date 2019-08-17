
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	// proxymethod "./proxyMethod"
	proxymethod "CLI/proxyMethod"

	"github.com/urfave/cli"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	isDaemon = false
)

var (
	filename = os.Getenv("HOME") + "/go/src/CLI/defaultConfig.json"
)

func main() {
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
			Name:   "run",
			Usage:  "Run the proxy server",
			Action: actionRun,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "daemon,d",
					Usage:       "daemon flag",
					Destination: &isDaemon,
				},
			},
		},
		{
			Name:  "reload",
			Usage: "Reload the proxy server. Stops previous one and opens another with new configuration file",
			Action: func(c *cli.Context) error {
				actionRun(c)
				actionStop(c)
				fmt.Println("Server is reloaded")
				return nil
			},
		},
		{
			Name:   "stop",
			Usage:  "Completely stops the server",
			Action: actionStop,
		},
	}

	err := app.Run(os.Args)
	Check(err)

}

func actionStop(c *cli.Context) error {
	_, err := os.Stat(getPidFilePath())
	if err != nil {
		log.Println("File is not running")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(getPidFilePath())
	if err != nil {
		log.Println("Unable to read file")
		os.Exit(1)
	}

	ProcessID, err := strconv.Atoi(string(data))
	if err != nil {
		log.Println("Unable to parse")
		os.Exit(1)
	}

	process, err := os.FindProcess(ProcessID)
	if err != nil {
		log.Println("Unable to parse")
		os.Exit(1)
	}

	os.Remove(getPidFilePath())

	fmt.Printf("Stop PID: %v", ProcessID)

	err = process.Kill()
	if err != nil {
		log.Println("Unable to kill process")
		os.Exit(1)
	}

	fmt.Printf("PID %v stopped!", ProcessID)
	return nil

}

func actionRun(c *cli.Context) error {
	if isDaemon {
		return runDaemon(c)
	}
	fmt.Println("Server is running\n")
	fmt.Printf("%v will be used as configuration file. If you changed your opinion, use reload command\n\n", filename)
	proxymethod.RunServer(filename)
	return nil
}

func runDaemon(ctx *cli.Context) error {
	if _, err := os.Stat(getPidFilePath()); err == nil {
		fmt.Println("Running..")
		os.Exit(1)
		return nil
	}
	cmd := exec.Command(os.Args[0], "run")
	cmd.Start()
	fmt.Println("Daemon process ID is : ", cmd.Process.Pid)
	savePID(cmd.Process.Pid)
	os.Exit(0)
	return nil
}

func getPidFilePath() string {
	return os.Getenv("HOME") + "/go/src/CLI/daemon.pid"
}

func savePID(pid int) {
	file, err := os.Create(getPidFilePath())
	if err != nil {
		log.Printf("Unable to create pid file:%v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))

	if err != nil {
		log.Printf("Unable to create pid file:%v\n", err)
		os.Exit(1)
	}
	file.Sync()
}
