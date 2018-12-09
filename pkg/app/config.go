package app

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

// Config type is the configuration entity of the application
type Config struct {
	Mode string
	Directories []string
}

// config file stuff
var (
	configFileName = "config"
	configFileExt = ".yml"
	configType = "yaml"
	appName = "gitbatch"

	configurationDirectory = filepath.Join(osConfigDirectory(), appName)
	configFileAbsPath = filepath.Join(configurationDirectory, configFileName)
)

// configuration items
var (
	modeKey = "mode"
	modeKeyDefault = "fetch"
	pathsKey = "paths"
	pathsKeyDefault = []string{"."}
)

// LoadConfiguration returns a Config struct is filled 
func LoadConfiguration() (*Config, error) {
	if err := initializeConfigurationManager(); err != nil {
		return nil, err
	}
	if err := setDefaults(); err != nil {
		return nil, err
	}
	if err := readConfiguration(); err != nil {
		return nil, err
	}
	config := &Config{
		Mode: viper.GetString(modeKey),
		Directories: viper.GetStringSlice(pathsKey),
	}
	return config, nil
}

// set default configuration parameters
func setDefaults() error {
	viper.SetDefault(modeKey, modeKeyDefault)
	// viper.SetDefault(pathsKey, pathsKeyDefault)
	return nil
}

// read configuration from file
func readConfiguration() error{
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		// if file does not exist, simply create one
		if _, err := os.Stat(configFileAbsPath+configFileExt); os.IsNotExist(err) {
			os.MkdirAll(configurationDirectory, 0755)
			os.Create(configFileAbsPath+configFileExt)
		} else {
			return err
		}
		// let's write defaults
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}

// write configuration to a file
func writeConfiguration() error{
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

// initialize the configuration manager
func initializeConfigurationManager() error {
	// config viper
	viper.AddConfigPath(configurationDirectory)
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configType)

	return nil
}

// returns OS dependent config directory
func osConfigDirectory() (osConfigDirectory string) {
	switch osname := runtime.GOOS; osname {
	case "windows":
		osConfigDirectory = os.Getenv("APPDATA")
	case "darwin":
		osConfigDirectory = os.Getenv("HOME") + "/Library/Application Support"
	case "linux":
		osConfigDirectory = os.Getenv("HOME") + "/.config"
	default:
		log.Warn("Operating system couldn't be recognized")
	}
	return osConfigDirectory
}