package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type ServiceConfig struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
	Port     int    `json:"port"`
}

const (
	COMMAND_POS   = 1
	FILE_PATH_POS = 2

	COMMAND_LENGTH = 2

	SUCCESS_EXIT_CODE = 0
	INVALID_EXIT_CODE = -1

	MAX_ALLOWED_PORT = 65535
)

func main() {
	if len(os.Args) < COMMAND_LENGTH {
		printUsage()
		os.Exit(INVALID_EXIT_CODE)
	}

	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	strictFlag := validateCmd.Bool("strict", false, "Enable strict validation")

	switch os.Args[COMMAND_POS] {
	case "validate":
		validateCmd.Parse(os.Args[FILE_PATH_POS:])

		if validateCmd.NArg() < 1 {
			fmt.Println("Error: Missing file path.")
			fmt.Println("Usage: servicectl validate <file.json>")
			os.Exit(INVALID_EXIT_CODE)
		}

		for _, file := range validateCmd.Args() {
			handleValidate(file, *strictFlag)
		}

	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(INVALID_EXIT_CODE)
	}
}

func handleValidate(filePath string, strict bool) {
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

func printUsage() {
	fmt.Println("Usage: servicectl <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  validate    Validate a service configuration file")
	fmt.Println("  help        Show this help message")
}
