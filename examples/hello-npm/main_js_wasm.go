package main

import (
	"log"

	"github.com/jcbhmr/goos-js2-go/syscall/js"
)

func await(jsPromise js.Value) js.Value {
	c := make(chan js.Value)
	jsHandleResolve := js.FuncOf(func(this js.Value, args []js.Value) any {
		c <- args[0]
		close(c)
		return js.Undefined()
	})
	defer jsHandleResolve.Release()
	jsHandleReject := js.FuncOf(func(this js.Value, args []js.Value) any {
		close(c)
		panic(js.Error{Value: args[0]})
	})
	defer jsHandleReject.Release()
	jsPromise.Call("then", jsHandleResolve, jsHandleReject)
	return <-c
}

var jsPrettier = await(js.Import("prettier", js.Undefined()))

func main() {
	code := `hello () ;;; world( true,  42)`
	formattedCode := await(jsPrettier.Call("format", code, map[string]any{
		"parser": "typescript",
	})).String()
	log.Printf("original code:\n%s", code)
	log.Println()
	log.Printf("formatted code:\n%s", formattedCode)
}
