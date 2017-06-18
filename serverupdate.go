package main

import "io"
import "os"
import "net/http"
import "archive/zip"
import "path/filepath"

import "github.com/codecat/go-libs/log"

const serverUpdateLocation = "tmp/ManiaplanetServer_Latest.zip"

func makeTempFolder() bool {
	//TODO: Use os.TempDir() instead
	if pathExists("tmp") {
		return true
	}

	err := os.Mkdir("tmp", os.ModePerm)
	if err != nil {
		log.Fatal("Couldn't create \"tmp\" folder: %s", err.Error())
		return false
	}

	return true
}

func downloadServerUpdate() bool {
	if !makeTempFolder() {
		return false
	}

	if pathExists(serverUpdateLocation) {
		log.Warn("Previous server download did not remove temporary file, removing now.")
		os.Remove(serverUpdateLocation)
	}

	out, err := os.Create(serverUpdateLocation)
	if err != nil {
		log.Fatal("Failed to create temporary download stream: %s", err.Error())
		return false
	}
	defer out.Close()

	resp, err := http.Get(Config.Server.LatestURL)
	if err != nil {
		log.Fatal("Failed to get from http: %s", err.Error())
		return false
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal("Failed to copy from http stream: %s", err.Error())
		return false
	}

	return true
}

func extractServerUpdate() bool {
	r, err := zip.OpenReader(serverUpdateLocation)
	if err != nil {
		log.Fatal("Couldn't open server zip file: %s", err.Error())
		return false
	}
	defer r.Close()

	if !pathExists("server") {
		os.Mkdir("server", os.ModePerm)
	}

	for _, f := range r.File {
		fi := f.FileInfo()

		if fi.IsDir() {
			log.Trace("Directory: \"%s\"", f.Name)
			if !pathExists("server/" + f.Name) {
				os.MkdirAll("server/" + f.Name, os.ModePerm)
			}
			continue
		}

		log.Trace("File: \"%s\"", f.Name)

		targetPath := "server/" + filepath.Dir(f.Name)
		if !pathExists(targetPath) {
			os.MkdirAll(targetPath, os.ModePerm)
		}

		ff, err := f.Open()
		if err != nil {
			log.Error("Couldn't load zip entry \"%s\": %s", f.Name, err.Error())
			continue
		}

		if pathExists("server/" + f.Name) {
			os.Remove("server/" + f.Name)
			log.Trace("Overwriting %s", f.Name)
		}

		out, err := os.Create("server/" + f.Name)
		out.Chmod(fi.Mode())
		if err != nil {
			log.Error("Couldn't create file \"%s\": %s", f.Name, err.Error())
			continue
		}

		io.Copy(out, ff)

		ff.Close()
		out.Close()
	}

	log.Info("Extracted %d files!", len(r.File))

	return true
}

func performServerUpdate() bool {
	if !downloadServerUpdate() {
		return false
	}

	log.Info("Successfully downloaded server update")

	stopAllServers()

	if !extractServerUpdate() {
		return false
	}

	os.Remove(serverUpdateLocation)
	return true
}

func serverUpdateCheck() {
	resp, err := http.Head(Config.Server.LatestURL)
	if err != nil {
		log.Fatal("Failed checking ManiaplanetServer version: %s", err.Error())
		return
	}

	lastModified := resp.Header.Get("Last-Modified")
	if lastModified == "" {
		log.Fatal("Failed checking Last-Modified of latest server zip")
		return
	}

	if lastModified == Config.Server.Version {
		log.Info("Server is up to date: \"%s\"", lastModified)
		return
	}

	defer saveConfig()

	log.Info("Upating Maniaplanet server to: \"%s\"", lastModified)
	if !performServerUpdate() {
		log.Fatal("Server update failed!")
	}
	Config.Server.Version = lastModified
}

func getPackUrl(id string) string {
	return "https://v4.live.maniaplanet.com/ingame/public/titles/download/" + id + ".Title.Pack.gbx"
}

func downloadPackUpdate(id string) bool {
	if !makeTempFolder() {
		return false
	}

	url := getPackUrl(id)
	target := "tmp/" + id + ".Title.Pack.Gbx"

	if pathExists(target) {
		log.Warn("Previous pack download did not move temporary file, removing now.")
		os.Remove(target)
	}

	out, err := os.Create(target)
	if err != nil {
		log.Fatal("Failed to create temporary download stream: %s", err.Error())
		return false
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to get from http: %s", err.Error())
		return false
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal("Failed to copy from http stream: %s", err.Error())
		return false
	}

	return true
}

func movePackUpdate(id string) bool {
	source := "tmp/" + id + ".Title.Pack.Gbx"
	target := "server/UserData/Packs/" + id + ".Title.Pack.Gbx"

	err := os.Rename(source, target)
	if err != nil {
		log.Fatal("Failed to move pack %s update: %s", id, err.Error())
		return false
	}
	return true
}

func performPackUpdate(id string) bool {
	if !downloadPackUpdate(id) {
		return false
	}

	log.Info("Successfully downloaded pack")

	stopAllServersWithPack(id)

	if !movePackUpdate(id) {
		return false
	}

	return true
}

func packUpdateCheck(id string) {
	if !pathExists("server") {
		log.Fatal("Tried checking for pack update while server is not downloaded!")
		return
	}

	url := getPackUrl(id)

	resp, err := http.Head(url)
	if err != nil {
		log.Fatal("Failed checking Maniaplanet pack %s version: %s", id, err.Error())
		return
	}

	lastModified := resp.Header.Get("Last-Modified")
	if lastModified == "" {
		log.Fatal("Failed checking Last-Modified of latest pack %s", id)
		return
	}

	needsToUpdate := true
	for _, v := range Config.Server.Packs {
		if v.ID != id {
			continue
		}
		if v.Version == lastModified {
			needsToUpdate = false
			break
		}
	}

	if !needsToUpdate {
		log.Info("Maniaplanet pack %s is up to date: \"%s\"", id, lastModified)
		return
	}

	defer saveConfig()

	log.Info("Updating Maniaplanet pack %s to: \"%s\"", id, lastModified)
	if !performPackUpdate(id) {
		log.Fatal("Maniaplanet pack update failed!")
	}

	for i := range Config.Server.Packs {
		if Config.Server.Packs[i].ID != id {
			continue
		}
		Config.Server.Packs[i].Version = lastModified
		return
	}

	Config.Server.Packs = append(Config.Server.Packs, ConfigPack{ id, lastModified })
}
