package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error obtaining config file path")
	}
	return homeDir + configFileName, err
}

func write(cfg Config) error {
	var marshalledData []byte
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting config file path")
	}
	marshalledData, err = json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshalling data")
	}
	os.WriteFile(configFilePath, marshalledData, 0600)
	return nil
}

func Read() (*Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return &Config{}, fmt.Errorf("error obtaining config file path")
	}
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return &Config{}, fmt.Errorf("error reading from file")
	}
	var config Config
	json.Unmarshal(content, &config)
	return &config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(*c)
	if err != nil {
		return fmt.Errorf("error setting user")
	}
	return nil
}
