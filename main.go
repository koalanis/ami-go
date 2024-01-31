// The main package is the entry point of the ami-go program and handles the CLI logic (for now)
package main

import (
	app "github.com/koalanis/ami-go/cmd"
)

func main() {
	app.Run()
}
