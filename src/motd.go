package src

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	defaultMotd = `┌─┐┌┬┐┌─┐┌┬┐┌┬┐┬ ┬
├┤ │││├─┘ │  │ └┬┘
└─┘┴ ┴┴   ┴  ┴  ┴   ` + version

	pathDynamicMotd = "/etc/emptty/motd-gen.sh"
	pathMotd        = "/etc/emptty/motd"
)

// Prints dynamic motd, if configured; otherwise prints motd, if pathMotd exists; otherwise it prints default motd.
func printMotd(conf *config) {
	if !conf.PrintMotd {
		return
	}
	if conf.DynamicMotd && fileIsExecutable(pathDynamicMotd) {
		cmd := exec.Command(pathDynamicMotd)
		dynamicMotd, err := cmd.Output()
		if err != nil {
			logPrint(err)
			printDefaultMotd()
			return
		}
		fmt.Print(revertColorEscaping(string(dynamicMotd)))
		if conf.DaemonMode {
			setColors(conf.FgColor, conf.BgColor)
		} else {
			resetColors()
		}
	} else if fileExists(pathMotd) {
		emptyMotd := true
		file, err := os.Open(pathMotd)
		defer file.Close()
		if err != nil {
			logPrint(err)
			printDefaultMotd()
			return
		}
		scan := bufio.NewScanner(file)
		for scan.Scan() {
			if scan.Text() != "" {
				emptyMotd = false
				fmt.Println(revertColorEscaping(scan.Text()))
			}
		}
		if !emptyMotd {
			if conf.DaemonMode {
				setColors(conf.FgColor, conf.BgColor)
			} else {
				resetColors()
			}
		}
	} else {
		printDefaultMotd()
	}
}

// Prints default motd.
func printDefaultMotd() {
	fmt.Printf("%s\n\n", defaultMotd)
}

// Reverts escaped color definitions to real color values.
func revertColorEscaping(value string) string {
	if value != "" {
		result := strings.ReplaceAll(value, "\\x1b", "\x1b")
		result = strings.ReplaceAll(result, "\\033", "\x1b")
		return result
	}
	return value
}

// Sets defined colors.
func setColors(fg, bg string) {
	color := ""

	if fg != "" {
		color += fg
	}
	if fg != "" && bg != "" {
		color += ";"
	}
	if bg != "" {
		color += bg
	}

	if fg == "" && bg == "" {
		color = "0"
	}
	fmt.Print("\x1b[0;" + color + "m\n")
}

// Resets colors to default.
func resetColors() {
	setColors("", "")
}
