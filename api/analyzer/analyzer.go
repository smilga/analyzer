package analyzer

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/smilga/analyzer/api"
)

var script = "puppeteer/analyze.js"

func Analyze(url string, service *api.Service) {
	s, err := json.Marshal(service)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	cmd := exec.Command("node", script, url, string(s))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
