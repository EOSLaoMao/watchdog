package eos

import (
	"time"

	"github.com/EOSLaoMao/watchdog/internal/engine"
	"github.com/gin-gonic/gin"
)

const (
	BlockProduceTime = 12 * 0.5
	ListenBlockPath  = "/eos/block/status"
	bpCount          = 21
	bpName           = "eoslaomaocom"
)

func Listen() {
	engine.E().GET(ListenBlockPath, func(c *gin.Context) {
		switch bs.Status {
		case BlockStatusPrepare:
			c.String(200, "preparing EOS monitor, <i>%v</i>", time.Now().Format(time.RFC1123))
		case BlockStatusOK:
			c.String(
				200,
				"%s in good condition :), current unpaid blocks is <b>%d</b>, current ranking: <b>%d</b>, <i>%v</i>",
				bpName,
				bs.UnpaidBlocks,
				bs.Ranking,
				time.Now().Format(time.RFC1123),
			)
		case BlockStatusDown:
			c.String(502, "<b>%s</b>, <i>%v</i>", string(bs.Status), time.Now().Format(time.RFC1123))
		}
	})
}
