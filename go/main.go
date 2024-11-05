package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/xuri/excelize/v2"
)

func main() {
	c := make(chan struct{})
	js.Global().Set("generator", js.FuncOf(generator))
	js.Global().Set("print", js.FuncOf(print))
	js.Global().Set("add", js.FuncOf(add))
	js.Global().Set("person", js.FuncOf(person))
	<-c
}

func generator(this js.Value, args []js.Value) interface{} {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	f.SetCellValue("Sheet1", "A1", args[0].String())

	dst := convertToUint8Array(f)

	return js.ValueOf(dst)
}

func convertToUint8Array(f *excelize.File) js.Value {
	buf := new(bytes.Buffer)
	f.Write(buf)

	src := buf.Bytes()
	dst := js.Global().Get("Uint8Array").New(len(src))
	js.CopyBytesToJS(dst, src)

	return dst
}

func print(this js.Value, args []js.Value) interface{} {
	arg1 := args[0]
	arg2 := args[1]

	return js.ValueOf(fmt.Sprintf("%s %s", arg1.String(), arg2.String()))
}

func add(this js.Value, args []js.Value) interface{} {
	result := args[0].Int() + args[1].Int()

	return js.ValueOf(result)
}

type Person struct {
	Name string
	Age  int
}

func (p *Person) Print() string {
	return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}

func person(this js.Value, args []js.Value) interface{} {
	p := &Person{
		Name: args[0].String(),
		Age:  args[1].Int(),
	}
	return js.ValueOf(p.Print())
}
