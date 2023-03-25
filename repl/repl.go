package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/yushyn-andriy/firefly/compiler"
	"github.com/yushyn-andriy/firefly/config"
	"github.com/yushyn-andriy/firefly/evaluator"
	"github.com/yushyn-andriy/firefly/lexer"
	"github.com/yushyn-andriy/firefly/object"
	"github.com/yushyn-andriy/firefly/parser"
	"github.com/yushyn-andriy/firefly/vm"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer, conf config.Config) {
	env := object.NewEnvironment()
	if conf.Mode == config.INTERACTIVE {
		scanner := bufio.NewScanner(in)
		for {
			fmt.Printf(PROMPT)
			scanned := scanner.Scan()
			if !scanned {
				return
			}

			line := scanner.Text()
			l := lexer.New(line)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParseErrors(out, p.Errors())
				continue
			}

			if conf.CompilerMode == true {
				comp := compiler.New()
				err := comp.Compile(program)
				if err != nil {
					fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
					continue
				}

				machine := vm.New(comp.Bytecode())
				err = machine.Run()
				if err != nil {
					fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
					continue
				}
				stackTop := machine.LastPoppedStackElem()
				if stackTop != nil {
					io.WriteString(out, stackTop.Inspect())
					io.WriteString(out, "\n")
				}
			} else {

				evaluated := evaluator.Eval(program, env)
				if evaluated != nil {
					io.WriteString(out, evaluated.Inspect())
					io.WriteString(out, "\n")
				}
			}
		}
	} else {
		input, err := ioutil.ReadAll(in)
		if err != nil {
			log.Fatal(err)
		}
		l := lexer.New(string(input))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
		}
		evaluator.Eval(program, env)
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
