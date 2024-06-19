package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ApiToken          string `json:"KANDJI_API_TOKEN"`
	ApiUrl            string `json:"KANDJI_API_URL"`
	StandardBlueprint string `json:"KANDJI_STANDARD_BLUEPRINT"`
	DevBlueprint      string `json:"KANDJI_DEV_BLUEPRINT"`
	TempPassword      string `json:"TEMP_PASSWORD"`
}

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	return &config, nil
}
