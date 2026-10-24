package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/service-operator/internal/config"
)

const (
	COMMAND_POS   = 1
	FILE_PATH_POS = 2

	COMMAND_LENGTH = 2

	INVALID_EXIT_CODE = -1
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
			config.HandleValidate(file, *strictFlag)
		}

	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(INVALID_EXIT_CODE)
	}
}

func printUsage() {
	fmt.Println("Usage: servicectl <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  validate    Validate a service configuration file")
	fmt.Println("  help        Show this help message")
}
