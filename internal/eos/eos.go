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
			c.String(200, "preparing EOS monitor, %v", time.Now().Format(time.RFC1123))
		case BlockStatusOK:
			c.String(
				200,
				"%s in good condition :), current unpaid blocks is %d, current ranking: %d, %v",
				bpName,
				bs.UnpaidBlocks,
				bs.Ranking,
				time.Now().Format(time.RFC1123),
			)
		case BlockStatusDown:
			c.String(502, "%s, %v", string(bs.Status), time.Now().Format(time.RFC1123))
		}
	})
}
