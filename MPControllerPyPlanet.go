package main

import "fmt"
import "bufio"
import "os"
import "os/exec"
import "time"

import "github.com/codecat/go-libs/log"

type PyPlanetController struct {
	Server *MPServer
	Configured bool

	Command *exec.Cmd
	KeepRunning bool
}

func (self *PyPlanetController) GetDirName() string {
	return fmt.Sprintf("mpman_%d", self.Server.ID())
}

func (self *PyPlanetController) GetPath() string {
	return "pyp/" + self.GetDirName() + "/"
}

func (self *PyPlanetController) GetDatabaseName() string {
	return fmt.Sprintf("%s_pyplanet_%d", Config.Database.Database, self.Server.ID())
}

func (self *PyPlanetController) EnsureDatabaseExists() bool {
	//TODO: Fuck MySQL.
	name := self.GetDatabaseName()

	res := dbQuery("SHOW DATABASES WHERE `Database`=?", name)
	if res == nil {
		return false
	}

	if len(res) == 0 {
		log.Info("Creating database for PyPlanet: \"%s\"", name)
		if !dbExec("CREATE DATABASE " + name) {
			return false
		}
		dbExec("ALTER SCHEMA " + name + " DEFAULT CHARACTER SET utf8mb4  DEFAULT COLLATE utf8mb4_unicode_ci")
		//TODO: GRANT SELECT, DELETE, INSERT, UPDATE  ON `mpman`.* TO 'mpman'@'localhost';
	}

	return true
}

func (self *PyPlanetController) InitProject() bool {
	cmdInit := exec.Command("pyplanet", "init_project", self.GetDirName())
	cmdInit.Dir = "pyp"
	err := cmdInit.Run()
	if err != nil {
		log.Fatal("PyPlanet init_project command failed: %s", err.Error())
		return false
	}
	return true
}

func (self *PyPlanetController) WriteConfig() bool {
	//TODO: Replace this config writing with simpler writing: https://github.com/PyPlanet/PyPlanet/issues/380
	basePath := self.GetPath() + "settings/base.py"
	if pathExists(basePath) {
		os.Remove(basePath)
	}

	out, err := os.Create(basePath)
	if err != nil {
		log.Fatal("Couldn't create PyPlanet settings base file \"%s\": %s", basePath, err.Error())
		return false
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	fmt.Fprintf(w, "import os\n")
	fmt.Fprintf(w, "ROOT_PATH = os.path.dirname(os.path.dirname(__file__))\n")
	fmt.Fprintf(w, "TMP_PATH = os.path.join(ROOT_PATH, 'tmp')\n")
	fmt.Fprintf(w, "if not os.path.exists(TMP_PATH):\n")
	fmt.Fprintf(w, "  os.mkdir(TMP_PATH)\n")
	fmt.Fprintf(w, "DEBUG = False\n")

	fmt.Fprintf(w, "POOLS = [ 'default' ]\n")

	superAdmins := dbQuery("SELECT * FROM superadmins WHERE Server=?", self.Server.ID())

	fmt.Fprintf(w, "OWNERS = { 'default': [\n")
	for _, a := range superAdmins {
		fmt.Fprintf(w, "  '%s',\n", a["Login"].(string))
	}
	fmt.Fprintf(w, "]}\n")

	//TODO: Fuck MySQL. Use some local database isntead.
	fmt.Fprintf(w, "DATABASES = { 'default': {\n")
	fmt.Fprintf(w, "  'ENGINE': 'peewee_async.MySQLDatabase',\n")
	fmt.Fprintf(w, "  'NAME': '%s',\n", self.GetDatabaseName())
	fmt.Fprintf(w, "  'OPTIONS': {\n")
	fmt.Fprintf(w, "    'host': '%s',\n", Config.Database.Hostname)
	fmt.Fprintf(w, "    'user': '%s',\n", Config.Database.Username)
	fmt.Fprintf(w, "    'password': '%s',\n", Config.Database.Password)
	fmt.Fprintf(w, "    'charset': 'utf8mb4'\n")
	fmt.Fprintf(w, "  }\n")
	fmt.Fprintf(w, "}}\n")

	fmt.Fprintf(w, "DEDICATED = { 'default': {\n")
	fmt.Fprintf(w, "  'HOST': '127.0.0.1',\n")
	fmt.Fprintf(w, "  'PORT': '%d',\n", self.Server.Info["PortRPC"].(int))
	fmt.Fprintf(w, "  'USER': 'SuperAdmin',\n")
	fmt.Fprintf(w, "  'PASSWORD': '%s'\n", self.Server.Passwords.SuperAdmin)
	fmt.Fprintf(w, "}}\n")

	fmt.Fprintf(w, "MAP_MATCHSETTINGS = {\n")
	fmt.Fprintf(w, "  'default': '%s'\n", self.Server.MatchSettingsFilename())
	fmt.Fprintf(w, "}\n")

	fmt.Fprintf(w, "STORAGE = { 'default': {\n")
	fmt.Fprintf(w, "  'DRIVER': 'pyplanet.core.storage.drivers.local.LocalDriver',\n")
	fmt.Fprintf(w, "  'OPTIONS': {}\n")
	fmt.Fprintf(w, "}}\n")

	return true
}

func (self *PyPlanetController) Configure(server *MPServer) bool {
	self.Server = server

	if !pathExists("pyp") {
		os.Mkdir("pyp", os.ModePerm)
	}

	if !pathExists(self.GetPath()) {
		if !self.InitProject() {
			return false
		}
	}

	if !self.WriteConfig() {
		return false
	}

	if !self.EnsureDatabaseExists() {
		return false
	}

	self.Configured = true
	return true
}

func (self *PyPlanetController) IsConfigured() bool {
	return self.Configured
}

func (self *PyPlanetController) Start() {
	self.KeepRunning = true
	self.Command = exec.Command("./manage.py", "start")
	//self.Command.Stdout = os.Stdout
	//self.Command.Stderr = os.Stderr
	self.Command.Dir = self.GetPath()
	go func(self *PyPlanetController) {
		//TODO: PyPlanet requires a bit of a delay after starting the MP server: https://github.com/PyPlanet/PyPlanet/issues/386
		// From my testing it takes Maniaplanet a lot of time to start the XMLRPC port
		time.Sleep(30000 * time.Millisecond)
		err := self.Command.Run()
		if err != nil {
			log.Warn("PyPlanet did not exit cleanly: %s", err.Error())
		}
		log.Info("PyPlanet exited")
		self.Command = nil
	}(self)
}

func (self *PyPlanetController) Stop() {
	self.KeepRunning = false
	//
}

func (self *PyPlanetController) IsRunning() bool {
	return self.Command != nil
}
