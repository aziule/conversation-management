package core

import (
	"io/ioutil"
	"encoding/json"
)

// Config is the struct that will hold the runtime configuration
type Config struct {
	ListeningPort int `json:"listening_port"`
	FbVerifyToken string `json:"fb_verify_token"`
	FbApiVersion string `json:"fb_api_version"`
	FbPageAccessToken string `json:"fb_page_access_token"`
}

// LoadConfig loads the configuration located at the given path
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var config Config

	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}