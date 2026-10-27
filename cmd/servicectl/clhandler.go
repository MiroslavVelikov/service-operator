package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/service-operator/internal/config"
)

type CommandType int

const (
	COMMAND_POS   = 1
	FILE_PATH_POS = 2

	COMMAND_LENGTH = 2

	INVALID_EXIT_CODE    = -1
	OCCURRED_FAILER_CODE = 1
	COMMAND_VALIDATE     = "validate"
	COMMAND_CHECK        = "check"
	COMMAND_HELP         = "help"
)

func main() {
	if len(os.Args) < COMMAND_LENGTH {
		printUsage()
		os.Exit(INVALID_EXIT_CODE)
	}

	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	strictFlag := validateCmd.Bool("strict", false, "Enable strict validation")

	switch os.Args[COMMAND_POS] {
	case COMMAND_VALIDATE:
		validateCmd.Parse(os.Args[FILE_PATH_POS:])

		if validateCmd.NArg() < 1 {
			fmt.Println("Error: Missing file path.")
			fmt.Println("Usage: servicectl validate <file.json>")
			os.Exit(INVALID_EXIT_CODE)
		}

		for _, file := range validateCmd.Args() {
			config.HandleValidate(file, *strictFlag)
		}
	case COMMAND_CHECK:
		validateCmd.Parse(os.Args[FILE_PATH_POS:])

		if validateCmd.NArg() < 1 {
			fmt.Println("Error: Missing URL paths.")
			fmt.Println("Usage: servicectl check <URL...>")
			os.Exit(INVALID_EXIT_CODE)
		}

		hC := http.Client{Timeout: time.Duration(10) * time.Second}
		var wg sync.WaitGroup
		var failed atomic.Bool
		for _, path := range validateCmd.Args() {
			wg.Go(func() {
				now := time.Now()
				resp, err := hC.Get(path)
				duration := time.Since(now)
				if err != nil {
					fmt.Printf("%s, failed, %v, %v\n", path, duration, err)
					failed.Store(true)
					return
				}
				defer resp.Body.Close()

				errMsg := "no errors occurred"
				if resp.StatusCode/100 != 2 {
					failed.Store(true)
					if resp.StatusCode/100 == 5 {
						errMsg = fmt.Sprintf("HTTP failure: %s", resp.Status)
					}
				}
				fmt.Printf("%s, %d, %v, %s\n", path, resp.StatusCode, duration, errMsg)
			})
		}
		wg.Wait()
		if failed.Load() {
			os.Exit(OCCURRED_FAILER_CODE)
		}
	case COMMAND_HELP:
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
	fmt.Println("  check       Check given URL and return: URL, status code, duration, error")
	fmt.Println("  help        Show this help message")
}

type CheckOutput struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Duration   int    `json:"duration"`
	Error      string `json:"error"`
}
