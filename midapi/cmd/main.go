package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	apiconf "wudian_go/midapi/config"
	"wudian_go/midapi/handlers"
	apiserver "wudian_go/midapi/server"
	"wudian_go/lib/log"

	"go.uber.org/zap"
	"gopkg.in/urfave/cli.v1"
)

var (
	VersionCommand = cli.Command{
		Name:   "version",
		Action: apiconf.ShowVersion,
		Usage:  "print version and compile time info",
	}
	logger *zap.Logger
)

func interceptSignal(app *cli.App) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		app.RunAndExitOnError()
	}()
}

func test() {
	a := map[string]interface{}{
		"a": 1,
	}
	var b interface{}
	b = a
	c, _ := b.(map[string]interface{})
	fmt.Println(c["a"])

}

func main() {
	// test()
	app := cli.NewApp()
	app.Name = "node midapi"
	app.Usage = "node midapi service"
	app.Version = apiconf.Version()
	app.Action = func(ctx *cli.Context) error {

		fmt.Println("----------------------------------------")
		apiconf.Init(ctx.String("config"))
		logger = log.Initialize("", apiconf.Env(), path.Join(apiconf.LogPath(), "api.output.log"), path.Join(apiconf.LogPath(), "api.err.log"))

		interceptSignal(app)

		server := apiserver.NewServer()
		server.Start()

		return nil
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "./config.json",
			Usage: "config file",
		},
	}

	app.Commands = []cli.Command{
		VersionCommand,
	}

	handlers.InitLogger(logger)
	app.Run(os.Args)
}
