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
var Null = js.Null
var Undefined = js.Undefined
var ValueOf = js.ValueOf

const TypeBoolean = js.TypeBoolean
const TypeFunction = js.TypeFunction
const TypeNull = js.TypeNull
const TypeNumber = js.TypeNumber
const TypeObject = js.TypeObject
const TypeString = js.TypeString
const TypeSymbol = js.TypeSymbol
const TypeUndefined = js.TypeUndefined

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

var jsImportMeta = jsGo.Get("_importMeta")

func ImportMeta() Value {
	return jsImportMeta
}

var jsGlobalThis = js.Global().Get("globalThis")

func Global() Value {
	return jsGlobalThis
}
