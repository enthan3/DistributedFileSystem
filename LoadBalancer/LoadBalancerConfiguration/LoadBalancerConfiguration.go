package LoadBalancerConfiguration

import (
	"encoding/json"
	"os"
)

type LoadBalancerConfiguration struct {
	CurrentHTTPAddress string            `json:"LoadBalancerHTTPAddress"`
	CurrentRPCAddress  string            `json:"LoadBalancerRPCAddress"`
	MastersAddress     map[string]string `json:"LoadBalancerMasterAddress"`
	ServiceHTTPAddress string            `json:"FrontendServiceHTTPAddress"`
	ServiceRPCAddress  string            `json:"FrontendServiceRPCAddress"`
	HeartbeatDuration  int               `json:"LoadBalancerHeartbeatDuration"`
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
