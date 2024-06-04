package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/adrg/xdg"
	os2 "github.com/jcbhmr/goos-js2-go/internal/os"
	exec "golang.org/x/sys/execabs"
)

//go:embed go.patch
var goPatch []byte

func CreateCustomizedGOROOT() (string, error) {
	customizedGOROOT, err := xdg.DataFile("goos-js2-go/customized-goroot")
	if err != nil {
		return "", err
	}
	log.Printf("customizedGOROOT=%#+v", customizedGOROOT)

	err = os.RemoveAll(customizedGOROOT)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(customizedGOROOT, 0755)
	if err != nil {
		return "", err
	}

	err = os2.CopyFS(customizedGOROOT, os.DirFS(runtime.GOROOT()))
	if err != nil {
		return "", err
	}

	goPatch2, err := xdg.DataFile("goos-js2-go/go.patch")
	if err != nil {
		return "", err
	}
	err = os.WriteFile(goPatch2, goPatch, 0644)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "apply", goPatch2)
	cmd.Dir = customizedGOROOT
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("cmd.Output() failed: %w\noutput:\n%s", err, out)
	}
	log.Printf("out=%#+v", out)

	return customizedGOROOT, nil
}
