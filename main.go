package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dop251/goja_nodejs/require"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
)

func main() {
	registry := new(require.Registry)
	data, err := ioutil.ReadFile("lib/mermaid.js")
	strData := string(data)
	vm := goja.New()
	this := registry.Enable(vm)
	console.Enable(vm)
	if this == nil {
		panic(this)
	}
	//program, err := goja.Compile("lib/example.js", strData, true)
	v, err := vm.RunScript("main", strData)
	if err != nil {
		panic(err)
	}
	//num := v.Export().(int64)
	fmt.Println(v)
}
