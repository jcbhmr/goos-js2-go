## Installation

```sh
go install github.com/jcbhmr/gojs/cmd/...@latest
```

## Usage

### Runtime environment

You can run the output JavaScript in any ECMAScript 2020 environment that
implements the Web Assembly JavaScript Interface. That includes Node.js, the
browser, Deno, Bun, and more.

The following **optional** features may enhance your Go JavaScript WebAssembly
app:

- Console
- Fetch Standard
- `node:fs`
- `node:process`
- `node:child_process`
- `node:net`

If you're in the browser and want to provide some polyfills for these features
you may do so using a pattern like this:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <script type="importmap">
      {
        "scopes": {
          "./app.js": {
            "node:fs": "./my-fs.js",
            "node:process": "./my-process.js"
          }
        }
      }
    </script>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>My Go WebAssembly App</title>
    <script type="module" src="app.js"></script>
  </head>
  <body>
    <p>The file contents should be printed below</p>
    <pre id="file-contents">This text should be replaced by Go code!</pre>
  </body>
</html>
```

### Call exposed functions

Go JavaScript WebAssembly programs are always linear; they still have a `main()`
and an exit code. However you _can_ call exported functions from JavaScript if
you use a hack like this:

```go
var jsSayHello = js.FuncOf(func(this js.Value, args []js.Value) any {
  fmt.Println("Hello from Go!")
  return js.Undefined()
})

var jsExit = js.FuncOf(func(this js.Value, args []js.Value) any {
  os.Exit(0)
  return js.Undefined()
})

func main() {
  importMetaURL, err := url.parse(js.ImportMeta().Get("url").String())
  jsDeferred := js.Global().Get(importMetaURL.Fragment)
  jsExports := js.ValueOf(map[string]any{
    "sayHello": jsSayHello,
    "exit": jsExit,
  })
  jsDeferred.Call("resolve", jsExports)
  select{} // Lazy wait forever.
}
```

```js
globalThis.abc123 = Promise.withResolvers();
// Notice how we are NOT awaiting the promise returned by `import()`?
import("./app.js#abc123").catch(globalThis.abc123.reject);
const lib = await globalThis.abc123.promise;
delete globalThis.abc123;
lib.sayHello();
lib.exit();
// Now the Go program has exited and we can't call into it anymore.
// 🛑 Error: Program has exited.
lib.sayHello();
```

## Development

This project is a mix of three projects that share the same root folder:

1. **An npm project** that uses Vite to compile a bunch of nice TypeScript into
   a single `gort0*` file
2. **A Go project** CLI wrapper around `go` that sets `GOOS=js`, sets
   `GOARCH=wasm`, and links the output WebAssembly with the `gort0*` runtime
3. **A Go package collection** that re-implements parts of the Go standard
   library to use JavaScript APIs
