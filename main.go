package main

import (
	"github.com/immid/tgmid/cmd"
	"github.com/immid/tgmid/pkg/base"
)

var Version = "dev"

func main() {
	base.Version = Version
	cmd.Execute()
}
