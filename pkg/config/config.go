package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func InitConfig(c *cli.Context) error {
	fmt.Println("init config")
	// 设置环境变量的前缀
	viper.SetEnvPrefix("KTTZAPP")
	// 允许环境变量重写配置文件中的值
	viper.AutomaticEnv()

	env := viper.GetString("ENV")
	if env == "" {
		viper.SetConfigFile("./configs/config.prod.toml")
	} else {
		viper.SetConfigFile("./configs/config.dev.toml")
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func GetString(key string) string {
	return viper.GetString(key)
}
