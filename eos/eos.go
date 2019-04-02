package eos

import (
	"fmt"
	"net/http"
)

const bpName string = "eoslaomaocom"

func MonitorBlocks() {
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

	http.ListenAndServe(":8080", nil)
}
