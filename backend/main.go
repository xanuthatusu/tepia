package main

import (
	"fmt"

	"github.com/xanuthatusu/tepia/pkg/server"
)

func main() {
	fmt.Println("howdy")
	server := server.New()
	server.Start()
}
