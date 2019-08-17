package eos

import (
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
			c.String(200, "preparing EOS monitor")
		case BlockStatusOK:
			if bs.Ranking > 21 {
				c.String(
					200,
					"OK, ranking: <b>%d</b>",
					bs.Ranking,
				)
			} else {
				c.String(
					200,
					"OK, unpaid blocks is <b>%d</b>, ranking: <b>%d</b>",
					bs.UnpaidBlocks,
					bs.Ranking,
				)
			}
		case BlockStatusTimeout:
			c.String(504, "request timeout", string(bs.Status))
		case BlockStatusDown:
			c.String(502, "<b>%s</b>", string(bs.Status))
		}
	})
}
