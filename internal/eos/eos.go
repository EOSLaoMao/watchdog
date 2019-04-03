package eos

import (
	"fmt"
	"net/http"
	"time"
)

const (
	BlockProduceTime = 12 * 0.5
	ListenPath       = "/eos/block/status"
	bpCount          = 21
	bpName           = "eoslaomaocom"
)

func Listen() {
	http.HandleFunc(ListenPath, func(w http.ResponseWriter, r *http.Request) {
		switch bs.Status {
		case StatusPrepare:
			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf("preparing EOS monitor %v", time.Now())))
		case StatusOK:
			w.WriteHeader(200)
			w.Write(
				[]byte(
					fmt.Sprintf(
						"%s in good condition :), current unpaid blocks is %d, %v",
						bpName, bs.UnpaidBlocks, time.Now(),
					),
				),
			)
		case StatusDown:
			w.WriteHeader(502)
			w.Write([]byte(bs.Status))
		}
	})
}
