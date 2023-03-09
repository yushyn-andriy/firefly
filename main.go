package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/yushyn-andriy/firefly/config"
	"github.com/yushyn-andriy/firefly/repl"
)

var (
	debug = flag.Bool("d", false, "debug mode")
	comp  = flag.Bool("c", false, "compiler mode")
)

func main() {
	flag.Parse()
	args := flag.Args()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	conf := config.Config{
		Debug:        *debug,
		Mode:         config.INTERACTIVE,
		CompilerMode: *comp,
	}

	if len(args) > 0 {
		conf.Mode = config.FROM_FILE

		filePath := args[0]
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("%s", err)
		}

		repl.Start(file, os.Stdout, conf)
	} else {
		fmt.Printf("Hello %s! This is Elephant programming language.\n", user.Username)
		repl.Start(os.Stdin, os.Stdout, conf)
	}
}
