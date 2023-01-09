//go:build ignore

package main

import (
	"log"
	"net/http"

	"github.com/vugu/vugu/devutil"
)

func main() {
	l := "127.0.0.1:8844"
	log.Printf("Serving WASM at %q", l)

	wc := devutil.NewWasmCompiler().SetBuildDir("./cmd/client/").SetGenerateDir("./cmd/client/")
	mux := devutil.NewMux()

	mux.Match(devutil.NoFileExt, devutil.DefaultAutoReloadIndex)
	mux.Exact("/main.wasm", devutil.NewMainWasmHandler(wc))
	mux.Exact("/wasm_exec.js", devutil.NewWasmExecJSHandler(wc))
	mux.Default(devutil.NewFileServer().SetDir("./cmd/client/"))

	log.Fatal(http.ListenAndServe(l, mux))
}
