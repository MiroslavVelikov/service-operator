package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	INVALID_EXIT_CODE = -1
	MAX_ALLOWED_PORT  = 65535
)

type ServiceConfig struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
	Port     int    `json:"port"`
}

func HandleValidate(filePath string, strict bool) {
	fmt.Printf("Validating %s (strict mode: %v)...\n", filePath, strict)

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(INVALID_EXIT_CODE)
	}

	var config ServiceConfig
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Validation failed! Invalid JSON: %v\n", err)
		os.Exit(INVALID_EXIT_CODE)
	}

	if config.Name == "" {
		fmt.Println("Validation failed! 'name' field is required.")
		os.Exit(INVALID_EXIT_CODE)
	}

	if config.Image == "" {
		fmt.Println("Validation failed! 'image' field is required.")
		os.Exit(INVALID_EXIT_CODE)
	}

	if config.Replicas < 0 {
		fmt.Println("Validation failed! 'replicas' field should be greater than or equal to 0.")
		os.Exit(INVALID_EXIT_CODE)
	}

	if config.Port < 1 || config.Port > MAX_ALLOWED_PORT {
		fmt.Println("Validation failed! 'port' field is required to be in the range of 1-65535.")
		os.Exit(INVALID_EXIT_CODE)
	}

	index := strings.LastIndex(config.Image, ":")
	if index == -1 {
		config.Image = config.Image + ":latest"
	}

	fmt.Println("Validation successful!")
	fmt.Printf("Service Name: %s, Port: %d\n", config.Name, config.Port)
}
