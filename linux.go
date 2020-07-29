package desktop

import (
	"errors"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
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

	out, err := exec.Command("wmctrl", "-lp").Output()
	if err != nil {
		return []string{}, err
	}

	lines := strings.Split(string(out), "\n")
	pids := []int{}
	for _, window := range lines {
		wmctrlColumns := strings.Split(window, " ")
		if len(wmctrlColumns) != 1 {
			pid64, err := strconv.ParseInt(wmctrlColumns[3], 10, 64)
			pid := int(pid64)
			if err != nil {
				return []string{}, err
			}
			var found bool
			for _, addedPid := range pids {
				if pid == addedPid {
					found = true
				}
			}
			if !found {
				pids = append(pids, pid)
			}
		}
	}

	apps := []string{}
	for _, pid := range pids {
		process, err := ps.FindProcess(pid)
		if err != nil {
			return []string{}, err
		}
		apps = append(apps, process.Executable())
	}
	return apps, nil
}

// LinuxQuitApp ... Quit a desktop application for a linux system
func LinuxQuitApp(name string) error {
	if err := checkLinuxOS(); err != nil {
		return err
	}
	cleanedName := strings.ReplaceAll(name, " ", "\\ ")
	err := exec.Command("pkill", "-x", cleanedName).Run()
	return err
}

func checkLinuxOS() error {
	if runtime.GOOS != "linux" {
		return errors.New("Wrong OS, only linux system supported")
	}
	return nil
}
