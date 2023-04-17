package main

import (
	"time"
)

func main() {

	go releasecheck("k3s-io", "k3s")
	go releasecheck("clanktron", "dummy")

	for {
		time.Sleep(time.Minute)
	}
	// go releasecheck("rancher", "rancher")
	// releasecheck("kubewarden", "kubewarden-controller")

}
