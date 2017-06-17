package main

import "io"
import "os"
import "net/http"
import "archive/zip"
import "path/filepath"

import "nimble/log"

const serverUpdateLocation = "tmp/ManiaplanetServer_Latest.zip"

func downloadServerUpdate() bool {
	err := os.Mkdir("tmp", os.ModePerm)
	if err != nil {
		log.Fatal("Couldn't create \"tmp\" folder: %s", err.Error())
		return false
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

	if pathExists("server") {
		log.Info("Deleting \"server\" directory")
		os.RemoveAll("server")
	}
	os.Mkdir("server", os.ModePerm)

	for _, f := range r.File {
		fi := f.FileInfo()

		if fi.IsDir() {
			log.Trace("Directory: \"%s\"", f.Name)
			os.MkdirAll("server/" + f.Name, os.ModePerm)
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

func deleteServerUpdateTmpFolder() {
	err := os.RemoveAll("tmp")
	if err != nil {
		log.Error("Failed to remove \"tmp\" directory: %s", err.Error())
	}
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

	deleteServerUpdateTmpFolder()

	return true
}

func serverUpdateCheck() {
	if pathExists("tmp") {
		log.Warn("The temporary \"tmp\" directory still exists, removing it!")
		deleteServerUpdateTmpFolder()
	}

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

	log.Info("Upating Maniaplanet server to: \"%s\"", lastModified)
	if !performServerUpdate() {
		log.Fatal("Server update failed!")
	}
	Config.Server.Version = lastModified
}
