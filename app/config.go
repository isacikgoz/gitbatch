package app

import (
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
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
	pathsKeyDefault     = []string{"."}
	logLevelKey         = "loglevel"
	logLevelKeyDefault  = "error"
	qucikKey            = "quick"
	qucikKeyDefault     = false
	recursionKey        = "recursion"
	recursionKeyDefault = 1
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
	var directories []string
	if len(viper.GetStringSlice(pathsKey)) <= 0 {
		d, _ := os.Getwd()
		directories = []string{d}
	} else {
		directories = viper.GetStringSlice(pathsKey)
	}
	config := &Config{
		Directories: directories,
		LogLevel:    viper.GetString(logLevelKey),
		Depth:       viper.GetInt(recursionKey),
		QuickMode:   viper.GetBool(qucikKey),
		Mode:        viper.GetString(modeKey),
	}
	return config, nil
}

// set default configuration parameters
func setDefaults() error {
	viper.SetDefault(logLevelKey, logLevelKeyDefault)
	viper.SetDefault(qucikKey, qucikKeyDefault)
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
			os.MkdirAll(configurationDirectory, 0755)
			os.Create(configFileAbsPath + configFileExt)
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
func writeConfiguration() error {
	err := viper.WriteConfig()
	return err
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
func osConfigDirectory(osname string) (osConfigDirectory string) {
	switch osname {
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
