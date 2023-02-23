//go:build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/vugu/vugu/devutil"
)

func main() {
	ip := flag.String("ip", "127.0.0.1", "ip on what to serve")
	flag.Parse()

	l := fmt.Sprintf("%s:8844", *ip)
	log.Printf("Serving WASM at %q", l)

	wc := devutil.NewWasmCompiler().SetBuildDir("./wasm").SetGenerateDir("./pkg/view")

	mux := devutil.NewMux()

	mux.Match(devutil.NoFileExt, indexHTML)
	mux.Exact("/main.wasm", devutil.NewMainWasmHandler(wc))
	mux.Exact("/wasm_exec.js", devutil.NewWasmExecJSHandler(wc))
	mux.Default(devutil.NewFileServer().SetDir("./wasm"))

	log.Fatal(http.ListenAndServe(l, mux))
}

var indexHTML = devutil.DefaultAutoReloadIndex.Replace(
	"<title>Vugu App</title>",
	"<title>SwimLogs</title>",
).Replace(
	"<!-- styles -->",
	`<meta name="viewport" content="width=device-width, initial-scale=1" />
	<link href="/tailwind.css" rel="stylesheet">`,
).Replace(
	"<!-- scripts -->",
	`<script src="https://kit.fontawesome.com/43a06af138.js" crossorigin="anonymous"></script>`)
