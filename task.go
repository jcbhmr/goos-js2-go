//usr/bin/true; exec go run "$0" "$@"
//go:build ignore

package main

import (
	"log"
	"os"
	"path/filepath"

	exec "golang.org/x/sys/execabs"
)

func Diff() error {
	goDir, err := filepath.Abs("go")
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "diff", "--binary")
	cmd.Dir = goDir
	diff, err := cmd.Output()
	if err != nil {
		return err
	}

	err = os.WriteFile("go.patch", diff, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetFlags(0)
	err := map[string]func() error{
		"diff": Diff,
	}[os.Args[1]]()
	if err != nil {
		log.Fatal(err)
	}
}
