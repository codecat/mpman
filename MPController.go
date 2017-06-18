package main

type MPController interface {
	Configure(server *MPServer) bool
	IsConfigured() bool

	Start()
	Stop()
	IsRunning() bool
}
