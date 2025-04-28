package configs

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var EnvConfig Config

type Config struct {
	Websocket Websocket `yaml:"websocket"`
}

type Websocket struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func init() {
	initViper()
}

func initViper() {
	vp := viper.New()

	vp.AddConfigPath("configs/")
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	vp.SetEnvKeyReplacer(replacer)

	err := vp.ReadInConfig()
	if err != nil {
		log.Fatalf("read config file error: %s", err)
	}
	err = vp.Unmarshal(&EnvConfig)
	if err != nil {
		log.Fatalf("unmarshal config file error: %s", err)
	}

}
