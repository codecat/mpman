package main

import "os"
import "time"
import "math/rand"

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

const genPasswordLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ 0123456789!@#$^*()"
func genPassword(n int) string {
	ret := make([]byte, n)
	for i := range ret {
		ret[i] = genPasswordLetters[rand.Intn(len(genPasswordLetters))]
	}
	return string(ret)
}

func seedRandom() {
	rand.Seed(time.Now().UTC().UnixNano())
}
