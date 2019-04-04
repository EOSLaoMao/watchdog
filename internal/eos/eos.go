package eos

import (
	"fmt"
	"net/http"
	"time"
)

const (
	BlockProduceTime = 12 * 0.5
	ListenBlockPath  = "/eos/block/status"
	bpCount          = 21
	bpName           = "eoslaomaocom"
)

func Listen() {
	http.HandleFunc(ListenBlockPath, func(w http.ResponseWriter, r *http.Request) {
		switch bs.Status {
		case BlockStatusPrepare:
			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf("preparing EOS monitor %v", time.Now())))
		case BlockStatusOK:
			w.WriteHeader(200)
			w.Write(
				[]byte(
					fmt.Sprintf(
						"%s in good condition :), current unpaid blocks is %d, %v",
						bpName, bs.UnpaidBlocks, time.Now(),
					),
				),
			)
		case BlockStatusDown:
			w.WriteHeader(502)
			w.Write([]byte(bs.Status))
		}
	})
}
