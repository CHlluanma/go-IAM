package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// CompletedConfig 用以表示配置已经完成
type CompletedConfig struct {
	*Config
}

func (c CompletedConfig) NewServer() (*GenericServer, error) {
	gin.SetMode(c.Mode)

	s := &GenericServer{
		Config:          c.Config,
		Engine:          gin.New(),
		ShutdownTimeout: 5,
	}

	initGenericServer(s)

	return s, nil
}

var (
	homeDir   = ".iam"
	envPrefix = "IAM"
)

// LoadConfig 读取配置文件和ENV变量
func LoadConfig(cfg string, defaultName string) {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath(".")
		dir, _ := os.UserHomeDir()
		viper.AddConfigPath(filepath.Join(dir, homeDir))
		viper.AddConfigPath("/etc/iam")
		viper.SetConfigName(defaultName)
	}

	// Use config file from the flag.
	viper.SetConfigType("yaml")   // set the type of the configuration to yaml.
	viper.AutomaticEnv()          // read in environment variables that match.
	viper.SetEnvPrefix(envPrefix) // set ENVIRONMENT variables prefix to IAM.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
