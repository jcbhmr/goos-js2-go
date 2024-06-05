package main

import (
	"time"

	"github.com/jcbhmr/goos-js2-go/syscall/js"
)

var jsDocument = js.Global().Get("document")
var jsCurrentTimePre = jsDocument.Call("querySelector", "#current-time")
var jsDateConstructor = js.Global().Get("Date")

func main() {
	for {
		currentTime := jsDateConstructor.New().Call("toLocaleTimeString").String()
		jsCurrentTimePre.Set("textContent", currentTime)
		time.Sleep(900 * time.Millisecond)
	}
}
