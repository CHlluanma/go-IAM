package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ahang7/go-IAM/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	configFlagName = "config"
	configFileType = "yaml"
)

var configFlagFile string
var configIn string

func init() {
	pflag.StringVarP(&configFlagFile, configFlagName, "c", configFlagFile, "set the configuration file, the default configuration file type is yaml")
}

func addConfigFile(prefixFlag string, configName string, fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(configFlagName))

	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(prefixFlag), "-", "_", -1))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	cobra.OnInitialize(func() {
		if configFlagFile != "" {
			viper.SetConfigFile(configFlagFile)
		} else {
			if configIn != "" {
				viper.AddConfigPath(configIn)
			} else {
				// 默认为当前包下的config包
				defaultIn := getRootDir()
				viper.AddConfigPath(defaultIn + "/config")
			}
			viper.SetConfigFile(configName)

		}
		viper.SetConfigType(configFileType)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("viper read config failed, err: %v", err)
			os.Exit(1)
		}
	})
}

func getRootDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("get Root Dir failed")
	}
	var infer func(dir string) string
	infer = func(dir string) string {
		modFile := filepath.Join(dir, "go.mod")
		if exist(modFile) {
			return dir
		}

		parent := filepath.Dir(dir)
		return infer(parent)
	}

	return infer(pwd)
}

func exist(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil || os.IsExist(err)
}

// SetConfigIn 设置配置文件路径
func SetConfigIn(in string) {
	configIn = in
}
