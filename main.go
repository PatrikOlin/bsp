package main

import (
	"encoding/gob"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

var hide bool
var fetch bool
var path = os.ExpandEnv("$HOME") + "/.bsp"

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
	nodeID := readFromFile()
	cmd := exec.Command("bspc", "node", "-g", "hidden")
	cmd.Run()

	focusCmd := exec.Command("bspc", "node", "-f", nodeID)
	focusCmd.Run()
}

func fetchNode() {
	saveCurrentNodeID()

	out, err := exec.Command("bspc", "query", "-N", "-n", ".hidden").Output()
	if err != nil {
		log.Fatalln(err)
	}
	nodeID := strings.TrimSpace(string(out))

	cmd := exec.Command("bspc", "node", nodeID, "-g", "hidden=off", "-f")
	cmd.Run()
}

func saveCurrentNodeID() {
	cmd, err := exec.Command("bspc", "query", "-N", "-n", ".focused").Output()
	if err != nil {
		log.Fatalln(err)
	}
	writeToFile(strings.TrimSpace(string(cmd)))
}

func writeToFile(nodeID string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalln("Error opening/creating file: ", err)
	}
	encoder := gob.NewEncoder(f)

	err = encoder.Encode(nodeID)
	if err != nil {
		log.Fatalln("Error encoding map: ", err)
	}

	f.Close()
}

func readFromFile() string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
		return ""

	}
	defer file.Close()

	var nodeID string
	decoder := gob.NewDecoder(file)
	decoder.Decode(&nodeID)

	return nodeID
}
