package SlaveConfiguration

import (
	"encoding/json"
	"os"
)

type SlaveConfiguration struct {
	CurrentRPCAddress string `json:"SlaveRPCAddress"`
	StoragePath       string `json:"SlaveStoragePath"`
}

func LoadSlaveConfiguration(path string) (SlaveConfiguration, error) {
	var config SlaveConfiguration
	configFile, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
