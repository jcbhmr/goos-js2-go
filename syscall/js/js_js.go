package js

import (
	"syscall/js"
	"unsafe"
)

type Error = js.Error
type Func = js.Func
type Type = js.Type
type Value = js.Value
type ValueError = js.ValueError
var CopyBytesToGo = js.CopyBytesToGo
var CopyBytesToJS = js.CopyBytesToJS
var FuncOf = js.FuncOf
var Global = js.Global
var Null = js.Null
var Undefined = js.Undefined
var ValueOf = js.ValueOf

var jsGo = func(id uint32, typeFlag byte) Value {
	type ref uint64
	const nanHead = 0x7FF80000
	v := struct {
		_     [0]func()
		ref   ref
		gcPtr *ref
	}{ref: (nanHead|ref(typeFlag))<<32 | ref(id)}
	return *(*Value)(unsafe.Pointer(&v))
}(6, 1)

func Import(specifier any, options any) Value {
	return jsGo.Call("_import", specifier, options)
}

var jsPromiseConstructor = Global().Get("Promise")

func Await(v Value) Value {
	jsPromise := jsPromiseConstructor.Call("resolve", v)
	c := make(chan Value)
	jsHandleResolve := FuncOf(func(this Value, args []Value) any {
		c <- args[0]
		return Value{}
	})
	jsHandleReject := FuncOf(func(this Value, args []Value) any {
		panic(&Error{Value: args[0]})
	})
	jsPromise.Call("then", jsHandleResolve, jsHandleReject)
	return <-c
}
