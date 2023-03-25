package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println(os.Getwd())
	err := http.ListenAndServe(":9090", http.FileServer(http.Dir("./webassembly/assets")))
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}
