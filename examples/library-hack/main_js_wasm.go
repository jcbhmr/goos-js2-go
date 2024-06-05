package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/jcbhmr/goos-js2-go/syscall/js"
)

var jsSayHello = js.FuncOf(func(this js.Value, args []js.Value) any {
	log.Println("Hello, WebAssembly!")
	return nil
})

func main() {
	log.SetFlags(0)
	importMetaURL := js.ImportMeta().Get("url").String()
	exports := map[string]any{
		"sayHello": jsSayHello,
	}
	js.Global().Get(strings.SplitAfter()[]).Call("resolve", exports)
	select {}
}
