package main

import "nimble/log"

var keepRunning = true

func main() {
	log.Open(log.CatTrace, log.CatFatal)
	log.Info("Initializing")

	loadConfigDefaults()

	if !pathExists("config.yaml") {
		log.Info("No config.yaml found - using empty config")
	} else if !loadConfig() {
		return
	}

	log.Config.MinLevel = Config.MinLogLevel
	log.Config.MaxLevel = Config.MaxLogLevel

	dbOpen()

	log.Info("Performing initial server update check")
	serverUpdateCheck()

	serverMonitor()

	log.Info("Shutting down")
	saveConfig()
}
