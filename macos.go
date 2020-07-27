package runningapps

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

// MacOSApplications ... Get a list of all running desktop application for a mac system
func MacOSApplications() ([]string, error) {
	if err := checkMacOS(); err != nil {
		return []string{}, err
	}
	out, err := exec.Command("osascript", "-e", `tell application "System Events" to get name of (processes where background only is false)`).Output()
	if err != nil {
		err := exec.Command("osascript", "-e", `tell application "System Events" to activate`).Run()
		return []string{}, err
	}

	dirtyApplications := strings.Split(string(out), ",")
	cleanApplications := []string{}
	for _, app := range dirtyApplications {
		if strings.TrimSuffix(strings.TrimSpace(app), "\n") != "Finder" {
			cleanApplications = append(cleanApplications, strings.TrimSpace(app))
		}
	}
	return cleanApplications, nil
}

// MacOSQuitApp ... Quit a desktop application for a mac system
func MacOSQuitApp(name string) error {
	if err := checkMacOS(); err != nil {
		return err
	}
	cleanedName := strings.ReplaceAll(name, " ", "\\ ")
	err := exec.Command("pkill", "-x", cleanedName).Run()
	return err
}

func checkMacOS() error {
	if runtime.GOOS != "darwin" {
		return errors.New("Wrong OS, only darwin system supported")
	}
	return nil
}
