package app

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// config file stuff
var (
	configFileName = "config"
	configFileExt  = ".yml"
	configType     = "yaml"
	appName        = "gitbatch"

	configurationDirectory = filepath.Join(osConfigDirectory(runtime.GOOS), appName)
	configFileAbsPath      = filepath.Join(configurationDirectory, configFileName)
)

// configuration items
var (
	modeKey             = "mode"
	modeKeyDefault      = "fetch"
	pathsKey            = "paths"
	quickKey            = "quick"
	quickKeyDefault     = false
	recursionKey        = "recursion"
	recursionKeyDefault = 1
)

// loadConfiguration returns a Config struct is filled
func loadConfiguration() (*Config, error) {
	if err := initializeConfigurationManager(); err != nil {
		return nil, err
	}
	if err := setDefaults(); err != nil {
		return nil, err
	}
	if err := readConfiguration(); err != nil {
		return nil, err
	}
	var directories []string
	if len(viper.GetStringSlice(pathsKey)) <= 0 {
		d, _ := os.Getwd()
		directories = []string{d}
	} else {
		directories = viper.GetStringSlice(pathsKey)
	}
	config := &Config{
		Directories: directories,
		Depth:       viper.GetInt(recursionKey),
		QuickMode:   viper.GetBool(quickKey),
		Mode:        viper.GetString(modeKey),
	}
	return config, nil
}

// set default configuration parameters
func setDefaults() error {
	viper.SetDefault(quickKey, quickKeyDefault)
	viper.SetDefault(recursionKey, recursionKeyDefault)
	viper.SetDefault(modeKey, modeKeyDefault)
	// viper.SetDefault(pathsKey, pathsKeyDefault)
	return nil
}

// read configuration from file
func readConfiguration() error {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		// if file does not exist, simply create one
		if _, err := os.Stat(configFileAbsPath + configFileExt); os.IsNotExist(err) {
			if err = os.MkdirAll(configurationDirectory, 0755); err != nil {
				return err
			}
			f, err := os.Create(configFileAbsPath + configFileExt)
			if err != nil {
				return err
			}
			defer f.Close()
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

// initialize the configuration manager
func initializeConfigurationManager() error {
	// config viper
	viper.AddConfigPath(configurationDirectory)
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configType)

	return nil
}

// returns OS dependent config directory
func osConfigDirectory(osName string) (osConfigDirectory string) {
	switch osName {
	case "windows":
		osConfigDirectory = os.Getenv("APPDATA")
	case "darwin":
		osConfigDirectory = os.Getenv("HOME") + "/Library/Application Support"
	case "linux":
		osConfigDirectory = os.Getenv("HOME") + "/.config"
	}
	return osConfigDirectory
}
