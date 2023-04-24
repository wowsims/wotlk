package main

import (
	"github.com/wowsims/wotlk/cmd/wowsimcli/cmd"
	"github.com/wowsims/wotlk/sim"
)

func init() {
	sim.RegisterAll()
}

func main() {
	cmd.Execute()
}
