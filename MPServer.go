package main

import "fmt"
import "os"
import "os/exec"

import "github.com/codecat/go-libs/log"

type MPServer struct {
	Info Row
	Command *exec.Cmd
	Configured bool

	Passwords struct {
		SuperAdmin string
		Admin string
	}

	Controller MPController
}

func (self *MPServer) Configure() bool {
	log.Info("Configuring server \"%s\"", self.Name())

	packID := self.Info["Title"].(string)
	log.Info("Checking for updates to pack %s", packID)
	packUpdateCheck(packID)

	self.Passwords.SuperAdmin = genPassword(20)
	self.Passwords.Admin = genPassword(20)

	pathUserData := "server/UserData/"
	pathConfig := pathUserData + "Config/" + self.ConfigFilename()
	pathMatchSettings := pathUserData + "Maps/MatchSettings/" + self.MatchSettingsFilename()

	log.Trace("Config path: \"%s\"", pathConfig)
	log.Trace("Match settings path: \"%s\"", pathMatchSettings)

	if pathExists(pathConfig) {
		os.Remove(pathConfig)
		log.Trace("Overwriting config")
	}
	if !pathExists(pathMatchSettings) {
		log.Info("Writing match settings file")
		if !self.WriteMatchSettings(pathMatchSettings) {
			return false
		}
	}

	if !self.WriteConfig(pathConfig) {
		return false
	}

	switch self.Info["Controller"].(string) {
		case "pyplanet": self.Controller = new(PyPlanetController)
	}
	if self.Controller != nil {
		log.Info("Configuring server controller \"%s\"", self.Info["Controller"].(string))
		if !self.Controller.Configure(self) {
			log.Warn("Server controller could not be configured!")
		}
	}

	self.Configured = true
	return true
}

func (self *MPServer) Start() {
	if self.IsRunning() {
		log.Warn("Tried starting server \"%s\" which is already running!", self.Name())
		return
	}

	log.Info("Starting server \"%s\"", self.Name())

	args := []string{}
	args = append(args, fmt.Sprintf("/dedicated_cfg=%s", self.ConfigFilename()))
	args = append(args, fmt.Sprintf("/game_settings=MatchSettings/%s", self.MatchSettingsFilename()))
	args = append(args, "/nodaemon")

	self.Command = exec.Command("./ManiaPlanetServer", args...)
	self.Command.Dir = "server/"
	go func(self *MPServer) {
		err := self.Command.Run()
		if err != nil {
			log.Error("Error while running server command: %s", err.Error())
		}
		self.Command = nil
	}(self)

	if self.Controller != nil {
		self.Controller.Start()
	}
}

func (self *MPServer) Stop() {
	if !self.IsRunning() {
		log.Warn("Tried stopping server \"%s\" which is not running!", self.Name())
		return
	}

	//

	if self.Controller != nil {
		self.Controller.Stop()
	}
}

func (self *MPServer) ID() int {
	return self.Info["ID"].(int)
}

func (self *MPServer) Name() string {
	return self.Info["Name"].(string)
}

func (self *MPServer) IsRunning() bool {
	return self.Command != nil
}

func (self *MPServer) ConfigFilename() string {
	return fmt.Sprintf("mpman_%d.txt", self.ID())
}

func (self *MPServer) MatchSettingsFilename() string {
	return fmt.Sprintf("mpman_maps_%d.txt", self.ID())
}
