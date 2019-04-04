package eos

import (
	"fmt"
	"net/http"
	"time"
)

const (
	BlockProduceTime  = 12 * 0.5
	ListenBlockPath   = "/eos/block/status"
	ListenRankingPath = "/eos/bpranking"
	bpCount           = 21
	bpName            = "eoslaomaocom"
)

func Listen() {
	http.HandleFunc(ListenBlockPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("....")
		switch bs.Status {
		case BlockStatusPrepare:
			w.WriteHeader(204)
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
	http.HandleFunc(ListenRankingPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("xxxx")
		switch {
		case ranking == 0:
			w.WriteHeader(204)
			w.Write([]byte(fmt.Sprintf("getting %s ranking %v", bpName, time.Now())))
		case ranking > 3:
			w.WriteHeader(502)
			w.Write([]byte(fmt.Sprintf("the %s ranking is seriously declining to %d", bpName, ranking)))
		default:
			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf("current %s ranking is %d", bpName, ranking)))
		}
	})
}
