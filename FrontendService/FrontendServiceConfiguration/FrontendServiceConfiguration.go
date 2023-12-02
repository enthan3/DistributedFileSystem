package FrontendServiceConfiguration

import (
	"encoding/json"
	"os"
)

type FrontendServiceConfiguration struct {
	Address     string `json:"FrontendServiceAddress"`
	CacheSize   int    `json:"FrontendServiceCacheSize"`
	StoragePath string `json:"FrontendServiceStoragePath"`
}

func LoadConfiguration(path string) (FrontendServiceConfiguration, error) {
	var config FrontendServiceConfiguration
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
