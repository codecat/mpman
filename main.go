package main

import "nimble/log"

func main() {
	log.Open(log.CatTrace, log.CatFatal)
	log.Info("Initializing")

	loadConfigDefaults()

	if !pathExists("config.yaml") {
		log.Info("No config.yaml found - using empty config")
	} else if !loadConfig("config.yaml") {
		return
	}

	log.Info("Performing initial server update check")
	serverUpdateCheck()

	saveConfig("config.yaml")

	serverMonitor()

	log.Info("Shutting down")
	saveConfig("config.yaml")
}
