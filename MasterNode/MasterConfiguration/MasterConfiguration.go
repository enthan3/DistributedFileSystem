package MasterConfiguration

import (
	"encoding/json"
	"os"
)

type MasterConfiguration struct {
	CurrentRPCAddress      string   `json:"MasterRPCAddress"`
	SlavesRPCAddress       []string `json:"SlavesRPCAddress"`
	MasterBackupRPCAddress string   `json:"MasterBackupRPCAddress"`
	LogPath                string   `json:"MasterLogPath"`
	HeartbeatDuration      int      `json:"MasterHeartbeatDuration"`
}

type MasterBackupConfiguration struct {
	CurrentRPCAddress string `json:"MasterBackupRPCAddress"`
	MasterRPCAddress  string `json:"MasterRPCAddress"`
	LogPath           string `json:"MasterBackupLogPath"`
	HeartbeatDuration int    `json:"MasterBackupHeartbeatDuration"`
}

func LoadMasterConfiguration(path string) (MasterConfiguration, error) {
	var config MasterConfiguration
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

func LoadMasterBackupConfiguration(path string) (MasterBackupConfiguration, error) {
	var config MasterBackupConfiguration
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
