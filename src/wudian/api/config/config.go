package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/urfave/cli.v1"
)

const version = "0.6.0"

var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
	COMMIT_VER string

	conf Config
)

func Init(configFile string) {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
		return
	}

	if err = json.Unmarshal(buf, &conf); err != nil {
		panic(err)
		return
	}
}

type Config struct {
	ListenAddress      string
	BackendCallAddress string
	Public             bool
	LogPath            string
	Debug              bool
	appName            string
	TimeOut            int
	TiConnEndpoint     string
	TiConnKey          string
	TiConnSecret       string
	Redis              string
	Redisidls          float64
	Redistimeout       float64
	Expire             float64
}

func Redis() string {
	return conf.Redis
}

func Redisidls() float64 {
	return conf.Redisidls
}
func Redistimeout() float64 {
	return conf.Redistimeout
}
func Expire() float64 {
	return conf.Expire
}

func LogPath() string {
	if len(conf.LogPath) != 0 {
		return conf.LogPath
	}
	return "logs"
}

func Env() string {
	if conf.Public {
		return "production"
	}
	return "development"
}

func ListenAddress() string {
	return conf.ListenAddress
}

func TimeOut() int {
	return conf.TimeOut
}

func BackendCallAddress() string {
	return conf.BackendCallAddress
}

func TiConnEndpoint() string {
	return conf.TiConnEndpoint
}

func TiConnKey() string {
	return conf.TiConnKey
}

func TiConnSecret() string {
	return conf.TiConnSecret
}

func Version() string {
	return version
}

func InitAppName(app string) {
	if len(conf.appName) > 0 {
		return
	}
	conf.appName = app
}

func Debug() bool {
	return conf.Debug
}

func AppName() string {
	return conf.appName
}

func ChainName() string {
	return conf.appName + "_chain"
}

func Public() bool {
	return conf.Public
}

func ShowVersion(ctx *cli.Context) {
	fmt.Println("version:", VERSION)
	fmt.Println("commit version:", COMMIT_VER)
	fmt.Println("build time:", BUILD_TIME)
	fmt.Println("go version:", GO_VERSION)
}
