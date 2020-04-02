package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dop251/goja_nodejs/console"

	"github.com/dop251/goja_nodejs/require"

	"github.com/dop251/goja"
)

func main() {
	registry := new(require.Registry)
	data, err := ioutil.ReadFile("../../lib/mermaid.min.js")
	strData := string(data)
	vm := goja.New()
	this := registry.Enable(vm)
	console.Enable(vm)

	if this == nil {
		panic(this)
	}
	v, err := vm.RunString(strData)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}
