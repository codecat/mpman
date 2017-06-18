package main

import "os"
import "io/ioutil"

import "nimble/log"
import "gopkg.in/yaml.v2"

type ConfigPack struct {
	ID string
	Version string
}

var Config struct {
	MinLogLevel int
	MaxLogLevel int

	Server struct {
		LatestURL string //TODO: Remove this one
		Version string

		Packs []ConfigPack
	}

	Database struct {
		Hostname string
		Port int
		Username string
		Password string
		Database string
	}
}

func loadConfigDefaults() {
	Config.MinLogLevel = log.CatTrace
	Config.MaxLogLevel = log.CatFatal

	Config.Server.LatestURL = "http://files.v04.maniaplanet.com/server/ManiaplanetServer_Latest.zip"

	Config.Database.Hostname = "localhost"
	Config.Database.Port = 3306
	Config.Database.Username = "root"
	Config.Database.Database = "mpman"
}

func readConfigFile() []byte {
	configData, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Couldn't read config file: %s", err.Error())
		return nil
	}
	return []byte(configData)
}

func loadConfig() bool {
	configData := readConfigFile()
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

func saveConfig() bool {
	configData, err := yaml.Marshal(Config)
	if err != nil {
		log.Fatal("Couldn't marshal yaml data: %s", err.Error())
		return false
	}

	err = ioutil.WriteFile("config.yaml", configData, os.ModePerm)
	if err != nil {
		log.Fatal("Couldn't write to file: %s", err.Error())
		return false
	}

	log.Info("Config saved")
	return true
}
