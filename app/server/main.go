package main

import (
	"net/http"

	"github.com/EOSLaoMao/watchdog/eos"
)

func main() {
	eos.WatchBlocks()

	http.ListenAndServe(":8080", nil)
}
