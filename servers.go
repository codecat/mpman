package main

import "time"

var keepRunning = true

func startServer() {
	//
}

func stopServer() {
	//
}

func serverMonitor() {
	for keepRunning {
		//TODO: Check for offline servers and turn them on if necessary
		//TODO: Check for online servers and turn them off if necessary

		time.Sleep(1 * time.Millisecond)
	}
}

func stopAllServers() {
	//TODO: Cleanly shut down all servers
}
