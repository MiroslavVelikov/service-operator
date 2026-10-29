package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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

	WORKERS_COUNT = 3

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

		var workers sync.WaitGroup
		jobs := make(chan string)
		outputs := make(chan CheckOutput, WORKERS_COUNT)
		client := http.Client{Timeout: time.Duration(10) * time.Second}
		failed := atomic.Bool{}

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		for w := 1; w <= WORKERS_COUNT; w++ {
			workers.Go(func() {
				for url := range jobs {
					if ctx.Err() != nil {
						return
					}

					result := checkURL(ctx, &client, url)
					if result.StatusCode/100 != 2 {
						failed.Store(true)
					}

					select {
					case outputs <- result:
					case <-ctx.Done():
						return
					}
				}
			})
		}

		go func() {
			defer close(jobs)

			for _, url := range validateCmd.Args() {
				select {
				case <-ctx.Done():
					return
				case jobs <- url:
				}
			}
		}()

		go func() {
			workers.Wait()
			close(outputs)
		}()

		for output := range outputs {
			fmt.Printf("%s, %v, %s, %s\n", output.URL, output.StatusCode, output.Duration, output.Error)
		}

		if failed.Load() || ctx.Err() != nil {
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
	URL        string        `json:"url"`
	StatusCode int           `json:"status_code"`
	Duration   time.Duration `json:"duration"`
	Error      string        `json:"error"`
}

func checkURL(
	ctx context.Context,
	client *http.Client,
	url string,
) CheckOutput {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return CheckOutput{
			URL:        url,
			StatusCode: 0,
			Duration:   0,
			Error:      err.Error(),
		}
	}

	now := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(now)
	if err != nil {
		return CheckOutput{
			URL:        url,
			StatusCode: http.StatusInternalServerError,
			Duration:   duration,
			Error:      err.Error(),
		}
	}
	defer resp.Body.Close()

	errMsg := "no errors occurred"
	if resp.StatusCode/100 != 2 {
		errMsg = fmt.Sprintf("HTTP failure: %s", resp.Status)
	}

	return CheckOutput{
		URL:        url,
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      errMsg,
	}
}
