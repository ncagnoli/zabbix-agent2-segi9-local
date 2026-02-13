package main

import (
	"fmt"
	"log"
	"os"

	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/plugin/container"
)

func main() {
	// Configure logging
	logFile := os.Getenv("SEGI9_LOG_FILE")
	if logFile != "" {
		// Open the log file securely (0600 - user only)
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			// Fallback to stderr if file cannot be opened
			log.SetOutput(os.Stderr)
			log.Printf("Failed to open log file %s: %v", logFile, err)
		} else {
			defer f.Close()
			log.SetOutput(f)
		}
	} else {
		// Default to stderr which Zabbix Agent captures
		log.SetOutput(os.Stderr)
	}

	// Capture panics to log file for professional error handling
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC: %v", r)
			fmt.Fprintf(os.Stderr, "Internal Error (Panic): %v\nCheck logs for details.\n", r)
			os.Exit(1)
		}
	}()

	log.Printf("Starting Segi9 plugin. Args: %v", os.Args)

	// Check for manual mode
	if len(os.Args) > 1 && os.Args[1] == "--manual" {
		runManual()
		return
	}

	h, err := container.NewHandler(impl.Name())
	if err != nil {
		errMsg := fmt.Sprintf("failed to create plugin handler %s", err.Error())
		log.Println(errMsg)
		panic(errMsg)
	}
	impl.Logger = &h

	log.Println("Handler created, executing...")

	err = h.Execute()
	if err != nil {
		errMsg := fmt.Sprintf("failed to execute plugin handler %s", err.Error())
		log.Println(errMsg)
		panic(errMsg)
	}
	log.Println("Plugin execution finished")
}

func runManual() {
	log.Println("Running in MANUAL mode")

	// Default configuration for manual mode
	// Provide a dummy GlobalOptions to avoid panic if Configure expects it
	globalOptions := &plugin.GlobalOptions{
		Timeout: 10,
	}

	privateOptions := &Config{Timeout: 10}

	impl.Configure(globalOptions, privateOptions)

	// Args for manual mode: ./plugin --manual <url> <authType> <user> <pass>
	// os.Args[0] is binary name
	// os.Args[1] is --manual
	// os.Args[2:] are parameters passed to Export

	params := []string{}
	if len(os.Args) > 2 {
		params = os.Args[2:]
	}

	// Do NOT log params here to avoid leaking secrets.
	// Export() already logs masked params.

	res, err := impl.Export("segi9.http", params, nil)
	if err != nil {
		log.Printf("Manual mode error: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Manual mode success")
	fmt.Println(res)
}
