package main

import (
	"fmt"
	"os"
	"os/user"
	"staq/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Print("The StaQ Programming Language")
	fmt.Printf("Version 0.0.1\n")
	fmt.Printf("Welcome, %s!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
