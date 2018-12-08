package app

import (

)

type Config struct {
	Mode string
	Directories []string
}

func LoadConfiguration() (*Config, error) {
	config := &Config{
	}
	return config, nil
}