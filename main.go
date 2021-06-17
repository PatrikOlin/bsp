package main

import (
	"flag"
	"log"
	"os/exec"
	"strings"
)

var hide bool
var fetch bool

func main() {
	flag.BoolVar(&hide, "hide", false, "hide node")
	flag.BoolVar(&fetch, "fetch", false, "fetch node")

	flag.Parse()

	if hide {
		hideNode()
	} else if fetch {
		fetchNode()
	} else {
		_, err := exec.Command("bspc", "query", "-N", "-n", ".hidden").Output()
		if err != nil {
			hideNode()
		} else {
			fetchNode()
		}
	}

}

func hideNode() {
	cmd := exec.Command("bspc", "node", "-g", "hidden")
	cmd.Run()
}

func fetchNode() {
	out, err := exec.Command("bspc", "query", "-N", "-n", ".hidden").Output()
	if err != nil {
		log.Fatalln(err)
	}
	nid := strings.TrimSpace(string(out))
	cmd := exec.Command("bspc", "node", nid, "-g", "hidden=off", "-f")
	cmd.Run()
}
