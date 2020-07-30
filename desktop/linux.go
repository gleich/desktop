package desktop

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// LinuxApplications ... Get a list of running desktop application for a linux system
func LinuxApplications() ([]string, error) {
	if err := checkLinuxOS(); err != nil {
		return []string{}, err
	}

	_, err := exec.LookPath("wmctrl")
	if err != nil {
		return []string{}, errors.New("The wmctrl tool is required to get a list of applications")
	}

	out, err := exec.Command("wmctrl", "-l").Output()
	if err != nil {
		return []string{}, err
	}

	lines := strings.Split(string(out), "\n")
	apps := []string{}
	for _, window := range lines {
		wmctrlColumns := strings.Split(window, " ")
		if len(wmctrlColumns) != 1 {
			windowID := wmctrlColumns[0]
			xpropcmd, err := exec.Command("xprop", "-id", windowID, "WM_CLASS").Output()
			if err != nil {
				return []string{}, err
			}
			xpropcmdChunks := strings.Split(string(xpropcmd), " ")
			app := strings.Trim(xpropcmdChunks[len(xpropcmdChunks)-1], "\"\n")
			var found bool
			for _, addedApp := range apps {
				if addedApp == app {
					found = true
				}
			}
			if !found {
				apps = append(apps, app)
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
		fmt.Println("Xprop failed")
		return err
	}
	xpropcmdChunks := strings.Split(string(xpropcmd), " ")
	pid := strings.Trim(xpropcmdChunks[len(xpropcmdChunks)-1], "\n")
	fmt.Printf("%#v", pid)

	err = exec.Command("kill", pid).Run()
	if err != nil {
		fmt.Println("Kill failed")
	}
	return err
}

func checkLinuxOS() error {
	if runtime.GOOS != "linux" {
		return errors.New("Wrong OS, only linux system supported")
	}
	return nil
}
