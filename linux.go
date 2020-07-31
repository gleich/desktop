package desktop

import (
	"errors"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// LinuxApplications ... Get a list of running desktop application for a linux system
func LinuxApplications() (map[string]int, error) {
	if err := checkLinuxOS(); err != nil {
		return map[string]int{}, err
	}

	_, err := exec.LookPath("wmctrl")
	if err != nil {
		return map[string]int{}, errors.New("The wmctrl tool is required to get a list of applications")
	}

	out, err := exec.Command("wmctrl", "-lp").Output()
	if err != nil {
		return map[string]int{}, err
	}

	lines := strings.Split(string(out), "\n")
	apps := map[string]int{}
	for _, window := range lines {
		wmctrlColumns := strings.Split(window, " ")
		if len(wmctrlColumns) != 1 {
			windowID := wmctrlColumns[0]
			appPID, err := strconv.Atoi(wmctrlColumns[3])
			if err != nil {
				return map[string]int{}, err
			}
			xpropcmd, err := exec.Command("xprop", "-id", windowID, "WM_CLASS").Output()
			if err != nil {
				return map[string]int{}, err
			}
			xpropcmdChunks := strings.Split(string(xpropcmd), " ")
			app := strings.Trim(xpropcmdChunks[len(xpropcmdChunks)-1], "\"\n")
			var found bool
			for addedApp := range apps {
				if addedApp == app {
					found = true
				}
			}
			if !found {
				apps[app] = appPID
			}
		}
	}
	return apps, nil
}

// LinuxQuitApp ... Quit a desktop application for a linux system
func LinuxQuitApp(name string) error {
	if err := checkLinuxOS(); err != nil {
		return err
	}
	xpropcmd, err := exec.Command("xprop", "-name", name, "_NET_WM_PID").Output()
	if err != nil {
		return err
	}
	xpropcmdChunks := strings.Split(string(xpropcmd), " ")
	pid := strings.Trim(xpropcmdChunks[len(xpropcmdChunks)-1], "\n")

	err = exec.Command("kill", pid).Run()
	return err
}

func checkLinuxOS() error {
	if runtime.GOOS != "linux" {
		return errors.New("Wrong OS, only linux system supported")
	}
	return nil
}
