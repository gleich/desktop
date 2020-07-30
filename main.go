package main

import (
	"fmt"

	"github.com/Matt-Gleich/desktop/desktop"
)

func main() {
	err := desktop.LinuxQuitApp("Firefox")
	if err != nil {
		fmt.Println(err)
	}
}
