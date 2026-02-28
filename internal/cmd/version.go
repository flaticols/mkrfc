package cmd

import (
	"fmt"
	"runtime/debug"
)

func printVersion() {
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "" {
		fmt.Println(info.Main.Version)
	} else {
		fmt.Println("dev")
	}
}
