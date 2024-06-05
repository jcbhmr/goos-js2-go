package main

import (
	_ "embed"
	"log"
	"runtime/debug"
)

//go:generate go run github.com/tdewolff/minify/cmd/minify -o wasm_exec.min.js wasm_exec.js
//go:embed wasm_exec.min.js
var wasmExecJS string

func main() {
	log.SetFlags(0)

	bi, _ := debug.ReadBuildInfo()
	log.Printf("%#+v", bi)
}
