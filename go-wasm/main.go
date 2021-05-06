package main

import (
	"syscall/js"
)

func callFuncWrapper() js.Func {
	callFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		url := args[0].String()
		callback := args[len(args)-1:][0]
		go CallH2C(url, func(data string) {
			callback.Invoke(data)
		})
		return true
	})
	return callFunc
}

func main() {
	js.Global().Set("getH2C", callFuncWrapper())
	select {} // Code must not finish
}
