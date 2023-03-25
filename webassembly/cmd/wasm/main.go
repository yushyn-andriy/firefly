package main

import (
	"bytes"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/yushyn-andriy/firefly/config"
	"github.com/yushyn-andriy/firefly/evaluator"
	"github.com/yushyn-andriy/firefly/repl"
)

func evaluate(input string) (string, error) {

	conf := config.Config{
		Debug:        false,
		Mode:         config.FROM_FILE,
		CompilerMode: false,
	}

	in := strings.NewReader(input)
	out := bytes.NewBufferString("")

	evaluator.SetStd(out, out)
	repl.Start(in, out, conf)

	return out.String(), nil
}

func evaluatorWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from", r)
			}
		}()

		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		inputText := args[0].String()

		// fmt.Printf("input %s\n", inputText)

		res, err := evaluate(inputText)
		if err != nil {
			fmt.Printf("unable to convert to json %s\n", err)
			return err.Error()
		}
		return res
	})
	return jsonFunc
}

func main() {

	fmt.Println("Firefly programming language ready for input...")
	js.Global().Set("evaluate", evaluatorWrapper())
	<-make(chan bool)
}
