package runningapps

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/Matt-Gleich/statuser/v2"
)

// MacOSApplications ... Get a list of all running desktop application for a mac system
func MacOSApplications() []string {
	if runtime.GOOS != "darwin" {
		statuser.ErrorMsg("MacOSApplications function only supported on a darwin based system", 1)
	}

	out, err := exec.Command("osascript", "-e", `tell application "System Events" to get name of (processes where background only is false)`).Output()
	if err != nil {
		err := exec.Command("osascript", "-e", `tell application "System Events" to activate`).Run()
		if err != nil {
			statuser.Error("Failed to get running list of applications", err, 1)
		}
	}

	dirtyApplications := strings.Split(string(out), ",")
	cleanApplications := []string{}
	for _, app := range dirtyApplications {
		if strings.TrimSuffix(strings.TrimSpace(app), "\n") != "Finder" {
			cleanApplications = append(cleanApplications, strings.TrimSpace(app))
		}
	}

	return cleanApplications
}
