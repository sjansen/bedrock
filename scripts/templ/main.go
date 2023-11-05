package main

import (
	"fmt"
	"os"

	"github.com/a-h/templ/cmd/templ/generatecmd"
)

func main() {
	err := generatecmd.Run(os.Stdout, generatecmd.Arguments{
		IncludeVersion: true,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
