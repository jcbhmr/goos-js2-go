package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/base64"
	"log"
	"os"

	exec "golang.org/x/sys/execabs"
)

//go:generate npm run build
//go:embed gort0_js_wasm.min.js
var gort0_js_wasm_min_js string

func main() {
	log.SetFlags(0)

	cmd := exec.Command("go", "build", "-o", "app.wasm", "./examples/hello-world")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS=js", "GOARCH=wasm")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("$ %s", cmd)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() %s: %v", cmd, err)
	}

	app_wasm, err := os.ReadFile("app.wasm")
	if err != nil {
		log.Fatalf("os.ReadFile() %s: %v", "app.wasm", err)
	}

	app_wasm_gz := bytes.Buffer{}
	gzipWriter := gzip.NewWriter(&app_wasm_gz)
	_, err = gzipWriter.Write(app_wasm)
	if err != nil {
		log.Fatalf("gzipWriter.Write(): %v", err)
	}
	err = gzipWriter.Close()
	if err != nil {
		log.Fatalf("gzipWriter.Close(): %v", err)
	}

	app_wasm_gz_base64 := base64.StdEncoding.EncodeToString(app_wasm_gz.Bytes())

	prelude := `//usr/bin/true; exec ${JS:-js} "$0" "$@"` + "\n" + `const __APP_WASM_GZ_BASE64__="` + app_wasm_gz_base64 + `";` + "\n"

	app_js := prelude + gort0_js_wasm_min_js

	err = os.WriteFile("app.js", []byte(app_js), 0644)
	if err != nil {
		log.Fatalf("os.WriteFile() %s: %v", "app.js", err)
	}
}
