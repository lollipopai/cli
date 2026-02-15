package main

import "github.com/lollipopai/cli/internal/cli"

var version = "dev"

func main() {
	cli.Execute(version)
}
