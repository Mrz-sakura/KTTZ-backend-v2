package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"kttz-server/internal/game"
	"kttz-server/pkg/config"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

func main() {
	app := cli.NewApp()

	// base application info
	app.Name = "kttz server"
	app.Version = "0.0.1"

	// flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "./configs/config.toml",
			Usage: "load configuration from `FILE`",
		},
		cli.BoolFlag{
			Name:  "cpuprofile",
			Usage: "enable cpu profile",
		},
	}

	app.Action = serve
	app.Run(os.Args)
}

func serve(c *cli.Context) error {
	err := config.InitConfig(c)
	if err != nil {
		panic(err)
	}

	log.SetFormatter(&log.TextFormatter{DisableColors: false})
	if viper.GetBool("core.debug") {
		log.SetLevel(log.DebugLevel)
	}

	if c.Bool("cpuprofile") {
		filename := fmt.Sprintf("cpuprofile-%d.pprof", time.Now().Unix())
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() { defer wg.Done(); game.Startup() }() // 开启游戏服
	//go func() { defer wg.Done(); web.Startup() }()  // 开启web服务器

	wg.Wait()
	return nil
}
