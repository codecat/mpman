package main

import "github.com/codecat/go-libs/log"

var keepRunning = true

func main() {
	seedRandom()

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

	if !dbOpen() {
		log.Fatal("Could not connect to database")
		return
	}

	log.Info("Performing initial server update check")
	serverUpdateCheck()

	serverMonitor()

	log.Info("Shutting down")
	saveConfig()
}
