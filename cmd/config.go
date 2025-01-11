package cmd

import (
	"encoding/json"
	"os"
)

type Config struct {
	Mode     string
	Database struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
		Port     int    `json:"port"`
		Sslmode  string `json:"sslmode"`
	} `json:"database"`
}

func loadConfig(file string) (*Config, error) {
	config := &Config{}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
