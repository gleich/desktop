package desktop

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

// MacOSApplications ... Get a list of all running desktop applications for a mac system
func MacOSApplications(includeMenubarApps bool) ([]string, error) {
	if err := checkMacOS(); err != nil {
		return []string{}, err
	}

	var command string
	if includeMenubarApps {
		command = `
			set nonAppleApps to ""
			tell application "System Events"
				set allApps to get name of processes
				repeat with oneApp in allApps
					set oneAppId to get bundle identifier of the process named oneApp
					if oneAppId does not contain "com.apple" and oneAppId is not missing value then
						if nonAppleApps is not "" then
							set nonAppleApps to nonAppleApps & ", "
						end if
						set nonAppleApps to nonAppleApps & oneApp
					end if
				end repeat
			end tell
			return nonAppleApps
		`
	} else {
		command = `tell application "System Events" to get name of (processes where background only is false)`
	}

	out, err := exec.Command("osascript", "-e", command).Output()

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
