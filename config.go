package main

import "os"
import "io/ioutil"

import "nimble/log"
import "gopkg.in/yaml.v2"

var Config struct {
	Server struct {
		LatestURL string
		Version string
	}

	Database struct {
		Hostname string
		Username string
		Password string
		Database string
	}
}

func loadConfigDefaults() {
	Config.Server.LatestURL = "http://files.v04.maniaplanet.com/server/ManiaplanetServer_Latest.zip"
}

func readConfigFile(fnm string) []byte {
	configData, err := ioutil.ReadFile(fnm)
	if err != nil {
		log.Fatal("Couldn't read config file: %s", err.Error())
		return nil
	}
	return []byte(configData)
}

func loadConfig(fnm string) bool {
	configData := readConfigFile(fnm)
	if configData == nil {
		return false
	}

	err := yaml.Unmarshal(configData, &Config)
	if err != nil {
		log.Fatal("Couldn't unmarshal yaml data: %s", err.Error())
		return false
	}

	log.Info("Config loaded")
	return true
}

func saveConfig(fnm string) bool {
	configData, err := yaml.Marshal(Config)
	if err != nil {
		log.Fatal("Couldn't marshal yaml data: %s", err.Error())
		return false
	}

	err = ioutil.WriteFile(fnm, configData, os.ModePerm)
	if err != nil {
		log.Fatal("Couldn't write to file: %s", err.Error())
		return false
	}

	log.Info("Config saved")
	return true
}
