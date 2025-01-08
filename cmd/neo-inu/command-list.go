package main

import (
	"neo-inu/internal"
	"neo-inu/pkg"
)

var commandList = []pkg.Command{
	internal.NewPingCommand(),
	internal.NewYgoCommand(),
}
