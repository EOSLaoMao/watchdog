package main

import (
	"github.com/EOSLaoMao/watchdog/internal/engine"
	"github.com/EOSLaoMao/watchdog/internal/eos"
	"github.com/EOSLaoMao/watchdog/internal/mixin"
	"github.com/EOSLaoMao/watchdog/internal/iost"
)

func main() {
	go eos.CheckUnpaidBlocks()
	go eos.Listen()

	go mixin.Listen()

	go iost.Listen()

	r := engine.E()
	r.Run()
}
