package LoadBalancerConfiguration

import (
	"encoding/json"
	"os"
)

type LoadBalancerConfiguration struct {
	Address   string `json:"Address"`
	CacheSize int64  `json:"CacheSize"`
}

func LoadConfiguration(path string) (LoadBalancerConfiguration, error) {
	var config LoadBalancerConfiguration
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
