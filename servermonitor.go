package main

import "time"

import "nimble/log"

var runningServers []*MPServer

func setupServer(info Row) {
	newServer := new(MPServer)
	newServer.Info = info
	runningServers = append(runningServers, newServer)
	newServer.Configure()
}

func isServerKnown(id int) bool {
	for _, s := range runningServers {
		if s.ID() == id {
			return true
		}
	}
	return false
}

func serverMonitor() {
	for keepRunning {
		log.Trace("Checking servers")

		servers := dbQuery("SELECT * FROM servers")
		for _, s := range servers {
			if !isServerKnown(s["ID"].(int)) {
				setupServer(s)
			}
		}

		for _, s := range runningServers {
			if !s.IsRunning() && s.Configured {
				s.Start()
			}
		}

		time.Sleep(5000 * time.Millisecond)
	}
}

func stopAllServers() {
	//TODO: Cleanly shut down all servers
}

func stopAllServersWithPack(id string) {
	//TODO: Cleanly shut down all servers with the given pack id
}
