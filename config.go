package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/apex/log"
)

type Config struct {
	HeaderURI                  string `json:"header.uri"`
	HeaderUsername             string `json:"header.username"`
	HTTPHostPort               string `json:"http.hostport"`
	HTTPPath                   string `json:"http.path"`
	PasswdFilename             string `json:"htpasswd.filepath"`
	MaintenanceIntervalSeconds int    `json:"maintenance.intervalSecs"`
	TokenByteLength            int    `json:"token.lengthBytes"`
	TokenValidityHours         int    `json:"token.durationHours"`
}

func NewConfig() *Config {
	return &Config{}
}

func (config *Config) Read(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(content, &config)
		if err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				log.WithError(err).WithField("offset", jsonErr.Offset).Error("Syntax error reading config")
			} else if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
				log.WithError(err).WithField("offset", jsonErr.Offset).Error("Unmarshal error reading config")
			} else {
				log.WithError(err).Errorf("Unknown error reading config: %T", err)
			}
		} else {
			err = config.validate()
		}
	}
	return err
}

func (config *Config) validate() error {
	var err error
	if config.HeaderUsername == "" {
		err = errors.New("Configuration file requires header.username")
	} else if config.HeaderURI == "" {
		err = errors.New("Configuration file requires header.uri")
	} else if config.HTTPHostPort == "" {
		err = errors.New("Configuration file requires http.hostport")
	} else if config.HTTPPath == "" {
		err = errors.New("Configuration file requires http.path")
	} else if config.PasswdFilename == "" {
		err = errors.New("Configuration file requires htpasswd.filepath")
	} else if config.TokenByteLength <= 0 {
		err = errors.New("Configuration file requires positive token.lengthBytes")
	} else if config.TokenValidityHours <= 0 {
		err = errors.New("Configuration file requires positive token.durationHours")
	} else if config.MaintenanceIntervalSeconds <= 0 {
		err = errors.New("Configuration file requires positive maintenance.intervalSecs")
	}
	return err
}
