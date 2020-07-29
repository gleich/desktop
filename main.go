package main

import (
	"github.com/Matt-Gleich/go_template/runningapps"
)

func main() {
	apps, _ := runningapps.LinuxApplications()
	runningapps.LinuxQuitApp(apps[2])
}
