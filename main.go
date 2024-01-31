package main

import (
	"fmt"
	"os"
	"os/user"
	"staq/repl"
)

const VERSION = "0.0.1"
const CODENAME = "Happiny"

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Println("The StaQ Programming Language")
	fmt.Printf("Version %s - %s\n", VERSION, CODENAME)
	fmt.Printf("Welcome, %s!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
