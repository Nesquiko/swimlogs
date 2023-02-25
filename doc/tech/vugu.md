# Vugu WASM library

## Vugu isn't ready / I don't want to implement everything

After almost building whole WASM frontend, I stumbled on a problem. I wanted to
reuse training-session-form.vugu component in training-edit-page.vugu, but I
couldn't.

My problem was that the training-session-form.vugu is implemented in a
way that it must be used as dynamic refrence in Vugu html part and instantiated
in the Init of parent component. This caused that Init of the
training-session-form.vugu was called when props to it weren't ready. It had to
be done in this way, because the vgform.Select component only populates a pointer
to String, but I needed it to populate values in Training, so I did that in final
validation function, which must be called to validate and set the values on Training.

I think I gave it a good try, I mean, I nearly built the whole thing in Vugu,
but I don't want to implement components, I want to build a frontend. I think
that Vugu isn't ready for production yet, or maybe I'm not ready for it.
Nevertheless, I switched to SolidJs

## Full HTML Mode

To add title and meta tags ./pkg/views/pages/root.vugu can use <html> tags.
But when developing, this will not reflect, because in ./cmd/client/wasmserver.go
there is a line:

`mux.Match(devutil.NoFileExt, devutil.DefaultAutoReloadIndex.Replace...`

This line serves an index.html with wrong title, meta tags. So for now there is
a hacky replace to default tags with correct ones. Didn't come up with better
solution.
