package main

import (
	"net/http"

	"github.com/EOSLaoMao/watchdog/internal/eos"
)

func main() {
	go eos.CheckUnpaidBlocks()
	go eos.Listen()

	http.ListenAndServe(":8080", nil)
}
