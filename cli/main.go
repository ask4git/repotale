package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login":
		fmt.Println("login: not implemented yet")
	case "connect":
		fmt.Println("connect: not implemented yet")
	case "open":
		fmt.Println("open: not implemented yet")
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("usage: repotale <login|connect|open>")
}
