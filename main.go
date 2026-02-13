package main

import (
	"fmt"
	"log"
	"os"

	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/plugin/container"
)

func main() {
	f, err := os.OpenFile("/tmp/zabbix-agent2-segi9.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// Fallback to stderr if file cannot be opened
		log.SetOutput(os.Stderr)
		log.Printf("Failed to open log file: %v", err)
	} else {
		defer f.Close()
		log.SetOutput(f)
	}

	log.Printf("Starting Segi9 plugin. Args: %v", os.Args)

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
