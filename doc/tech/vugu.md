# Vugu WASM library

## Full HTML Mode

To add title and meta tags ./pkg/views/pages/root.vugu can use <html> tags.
But when developing, this will not reflect, because in ./cmd/client/wasmserver.go
there is a line:

`mux.Match(devutil.NoFileExt, devutil.DefaultAutoReloadIndex.Replace...`

This line serves an index.html with wrong title, meta tags. So for now there is
a hacky replace to default tags with correct ones. Didn't come up with better
solution.

<!-- TODO: how to serve the index.html, or how to do it with only root.vugu?  -->
