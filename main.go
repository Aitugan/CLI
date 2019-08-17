
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

	app.Commands = []cli.Command{

		{
			Name:   "run",
			Usage:  "Run the proxy server",
			Action: run,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "daemon,d",
					Usage:       "daemon flag to run CLI on background",
					Destination: &isDaemon,
				},
			},
		},
		{
			Name:  "reload",
			Usage: "Reload the proxy server. Stops previous one and opens another with new configuration file",
			Action: func(c *cli.Context) error {
				run(c)
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

	fmt.Printf("PID %v stopped!\n", ProcessID)
	return nil

}

func run(c *cli.Context) {
	if c.Args().First() != "" {
		filename = c.Args().Get(0)
	}
	actionRun(c)
}

func actionRun(c *cli.Context) error {
	fmt.Println("Server is running\n")
	fmt.Printf("%v %T will be used as configuration file. If you changed your opinion, use reload command\n\n", filename, filename)

	if isDaemon {
		return runDaemon(c)
	}
	proxymethod.RunServer(filename)
	return nil
}

func runDaemon(ctx *cli.Context) error {
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
