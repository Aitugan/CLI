// reliable request response
// если с один бэкенд в раунд-робине не отвечает за указанное время, перейти к другому
// если один бэкенд в эникасте не отвечает, перейти к другому, если все не отвечают, повторить через определенное время, если и тогда никто не отвечает, вывести ошибку

// ТЕСТЫ НАПИСАААТЬ

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// _"net/http/pprof"

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//** request-response on rabbitmq
// запрос который будет отправляться, новый прокси метод(rabbitmq, например) его запросы отправляются в рэббит и там обрабатываются. также сервис который считывает с рэббита и оттуда возвращает ответ

// добавить логи
// сервер принимает запрос и пишет какие запросы принял и куда он отправляет запрос
// логировать метод HTTPMethod, path, header(json) http запросов

// добавить идентификатор на запрос и логировать его вмесе

// golang uuid

// сначала логи, потом интеграция с методом rabbit

// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"os/exec"

// 	proxymethod "proxymethod"
// 	"strconv"

// 	"github.com/urfave/cli"
// )

// func Check(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// var (
// 	isDaemon = false
// 	filename = os.Getenv("HOME") + "/go/src/CLI/defaultConfig.json"
// )

// func main() {
// 	app := cli.NewApp()
// 	app.Name = "CLI"
// 	app.Usage = "Idk, just a cli"
// 	app.Description = "For support, contact ..email.."
// 	app.UsageText = ""
// 	app.Version = "1.0.0"
// 	app.Authors = []cli.Author{
// 		cli.Author{
// 			Name: "Aitugan Mir",
// 		},
// 	}

// 	app.Commands = []cli.Command{

// 		{
// 			Name:   "run",
// 			Usage:  "Run the proxy server",
// 			Action: run,
// 			Flags: []cli.Flag{
// 				cli.BoolFlag{
// 					Name:        "daemon,d",
// 					Usage:       "daemon flag to run CLI on background",
// 					Destination: &isDaemon,
// 				},
// 			},
// 		},
// 		{
// 			Name:  "reload",
// 			Usage: "Reload the proxy server. Stops previous one and opens another with new configuration file",
// 			Action: func(c *cli.Context) error {
// 				if c.Args().First() != "" {
// 					filename = c.Args().Get(0)
// 				}
// 				actionStop(c)
// 				fmt.Println("Server is reloaded")
// 				actionRun(c)
// 				return nil
// 			},
// 		},
// 		{
// 			Name:   "stop",
// 			Usage:  "Completely stops the server",
// 			Action: actionStop,
// 		},
// 	}

// 	err := app.Run(os.Args)
// 	Check(err)

// }

// func actionStop(c *cli.Context) error {
// 	_, err := os.Stat(getPidFilePath())
// 	if err != nil {
// 		log.Println("File is not running")
// 		os.Exit(1)
// 	}

// 	data, err := ioutil.ReadFile(getPidFilePath())
// 	if err != nil {
// 		log.Println("Unable to read file")
// 		os.Exit(1)
// 	}

// 	ProcessID, err := strconv.Atoi(string(data))
// 	if err != nil {
// 		log.Println("Unable to parse")
// 		os.Exit(1)
// 	}

// 	process, err := os.FindProcess(ProcessID)
// 	if err != nil {
// 		log.Println("Unable to parse")
// 		os.Exit(1)
// 	}

// 	os.Remove(getPidFilePath())

// 	fmt.Printf("Stop PID: %v", ProcessID)

// 	err = process.Kill()
// 	if err != nil {
// 		log.Println("Unable to kill process")
// 		os.Exit(1)
// 	}

// 	fmt.Printf("PID %v stopped!\n", ProcessID)
// 	return nil

// }

// func run(c *cli.Context) {
// 	if c.Args().First() != "" {
// 		filename = c.Args().Get(0)
// 	}
// 	actionRun(c)
// }

// func actionRun(c *cli.Context) error {
// 	fmt.Println("Server is running\n")
// 	fmt.Printf("%v will be used as configuration file. If you changed your opinion, use reload command\n\n", filename)

// 	if isDaemon {
// 		return runDaemon(c)
// 	}
// 	proxymethod.RunServer(filename)
// 	// RunServer(filename)
// 	// http.HandleFunc(json.NewDecoder(os.Open(filename)), proxymethod.RunServer() //filename)
// 	return nil
// }

// func runDaemon(ctx *cli.Context) error {
// 	cmd := exec.Command(os.Args[0], "run")
// 	cmd.Start()
// 	fmt.Println("Daemon process ID is : ", cmd.Process.Pid)
// 	savePID(cmd.Process.Pid)
// 	os.Exit(0)
// 	return nil
// }

// func getPidFilePath() string {
// 	return os.Getenv("HOME") + "/go/src/CLI/daemon.pid"
// }

// func savePID(pid int) {
// 	file, err := os.Create(getPidFilePath())
// 	if err != nil {
// 		log.Printf("Unable to create pid file:%v\n", err)
// 		os.Exit(1)
// 	}
// 	defer file.Close()

// 	_, err = file.WriteString(strconv.Itoa(pid))

// 	if err != nil {
// 		log.Printf("Unable to create pid file:%v\n", err)
// 		os.Exit(1)
// 	}
// 	file.Sync()
// }

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	proxymethod "Final/proxyMethod"

	"github.com/urfave/cli"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// Global variables to check if server should run daemonly or not, and default file name
var (
	isDaemon = false
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
				if c.Args().First() != "" {
					filename = c.Args().Get(0)
				}
				actionStop(c)
				fmt.Println("Server is reloaded")
				//run(c)
				actionRun(c)
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

// Stops server
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
		filename = c.Args().First() //Get(0)
	}
	actionRun(c)
}

// actionRun runs server
func actionRun(c *cli.Context) error {
	fmt.Println("Server is running\n")
	fmt.Printf("%v %T will be used as configuration file. If you changed your opinion, use reload command\n\n", filename, filename)

	if isDaemon {
		return runDaemon(c)
	}
	proxymethod.RunServer(filename)
	return nil
}

// runDaemon is a command that runs server daemonly
func runDaemon(ctx *cli.Context) error {
	cmd := exec.Command(os.Args[0], "run")
	cmd.Start()
	fmt.Println("Daemon process ID is : ", cmd.Process.Pid)
	savePID(cmd.Process.Pid)
	os.Exit(0)
	return nil
}

// returns Pid filepath
func getPidFilePath() string {
	return os.Getenv("HOME") + "/go/src/CLI/daemon.pid"
}

// savePID creates and saves pid file
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
