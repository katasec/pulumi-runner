package main

import (
	"fmt"
)

func main() {
	fmt.Println("test")

	// Create a new pulumi program
	//p, err := createInlineProgram()

	p, err := createLocalProgram("hello", "/Users/ameerdeen/.ark/registry/ghcr.io/katasec/ark-resource-phello/v0.0.1")

	if err != nil {
		fmt.Printf("Error : %+v\n", err)
	} else {
		// Run Pulumi Up
		p.Preview()
	}
}
