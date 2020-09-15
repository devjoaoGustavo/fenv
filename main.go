package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devjoaoGustavo/fenv/service"
)

func main() {
	params, err := service.SearchParameters(os.Args[1:])
	if err != nil {
		log.Fatalf("fenv: %q", err)
	}
	result, err := service.GetParameters(params)
	if err != nil {
		log.Fatalf("fenv: %q", err)
	}
	for _, param := range result {
		fmt.Println(param)
	}
}
