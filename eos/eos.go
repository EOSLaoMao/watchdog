package eos

import (
	"fmt"
	"net/http"
)

const (
	BlockProduceTime = 12 * 0.5
	bpCount          = 21
	bpName           = "eoslaomaocom"
)

func WatchBlocks() {
	go func() {
		go logBlocks()
		go checkBlocks()

		c := make(chan struct{})
		<-c
	}()

	serveBlocks()
}

func serveBlocks() {
	http.HandleFunc("/block/status", func(w http.ResponseWriter, r *http.Request) {
		if s.Err != nil {
			w.WriteHeader(500)
			w.Write([]byte(s.Err.Error()))
		} else {
			w.Write(
				[]byte(
					fmt.Sprintf(
						"%s in good condition :), current unpaid blocks is %d, %v",
						bpName, s.UnpaidBlocks, s.Time,
					),
				),
			)
		}
	})
}
