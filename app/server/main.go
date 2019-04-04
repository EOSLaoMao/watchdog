package main

import (
	"github.com/EOSLaoMao/watchdog/internal/engine"
	"github.com/EOSLaoMao/watchdog/internal/eos"
)

func main() {
	go eos.CheckUnpaidBlocks()
	go eos.Listen()

	r := engine.E()
	r.Run()
}
