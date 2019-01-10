package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
)

var (
	// ErrInvalidConfig is returned when the configuration values fail to pass the criteria in the valid() function
	ErrInvalidConfig = errors.New("The configuration is invalid (see valid func)")
)

type WebCrawlerConfig struct {
	HTTPTimeout time.Duration `json:"http_timeout_seconds"`
	WorkerCount int           `json:"worker_count_integer"`
}

// LoadJSONConfig receives a file path for a JSON file, loads it onto the internal
// struct and returns that.
func LoadJSONConfig(fp string) (*WebCrawlerConfig, error) {
	fileBytes, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	var wcc WebCrawlerConfig
	err = json.Unmarshal(fileBytes, &wcc)
	if err != nil {
		return nil, err
	}

	// time.Duration needs to be asserted to seconds as I think
	// json.Unmarshal tries to be clever about how to parse the int in the JSON
	// to time.Duration
	wcc.HTTPTimeout = wcc.HTTPTimeout * time.Second

	if !wcc.Valid() {
		return nil, ErrInvalidConfig
	}

	return &wcc, nil
}

// Valid returns true if the configuration's values are OK, else it returns false.
func (wcc *WebCrawlerConfig) Valid() bool {
	return (wcc.HTTPTimeout > 0) && wcc.WorkerCount > 0
}
