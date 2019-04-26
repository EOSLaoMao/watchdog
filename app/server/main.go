package main

import (
	"github.com/EOSLaoMao/watchdog/internal/engine"
	"github.com/EOSLaoMao/watchdog/internal/eos"
	"github.com/EOSLaoMao/watchdog/internal/iotex"
	"github.com/EOSLaoMao/watchdog/internal/mixin"
)

func main() {
	go eos.CheckUnpaidBlocks()
	go eos.Listen()

	go mixin.Listen()

	go iotex.Listen()

	r := engine.E()
	r.Run()
}
