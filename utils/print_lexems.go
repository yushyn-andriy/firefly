package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/yushyn-andriy/firefly/lexer"
	"github.com/yushyn-andriy/firefly/token"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("Usage: print_lexems <path>")
	}

	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	l := lexer.New(string(data))
	for i := 0; ; i++ {
		tok := l.NextToken()
		fmt.Printf("%d: Type: %q, Literal: %q\n", i, tok.Type, tok.Literal)
		if tok.Type == token.EOF {
			break
		}
	}
}
