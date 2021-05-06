GOOS=js GOARCH=wasm go build -o  ../public/main.wasm
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js"  ../public/wasm_exec.js