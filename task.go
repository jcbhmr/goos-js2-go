//usr/bin/true; exec go run "$0" "$@"
//go:build ignore

package main

import (
	"log"
	"os"

	exec "golang.org/x/sys/execabs"
)

func Setup() error {
	cmd := exec.Command("npm", "install")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("$ %s", cmd)
	return cmd.Run()
}

func main() {
	log.SetFlags(0)
	task := map[string]func() error{
		"setup": Setup,
	}[os.Args[1]]
	if err := task(); err != nil {
		log.Fatal(err)
	}
}
