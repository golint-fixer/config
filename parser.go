package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type argument struct {
	Description string      `json:"description"`
	Default     interface{} `json:"default"`
	EnvName     string      `json:"env_name"`
	FlagName    string      `json:"flag_name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
}

func parseJSON(filename string) (map[string]argument, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c map[string]argument
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	ci := map[string]argument{}
	for k, v := range c {
		ci[strings.ToLower(k)] = v
	}

	return ci, nil
}
