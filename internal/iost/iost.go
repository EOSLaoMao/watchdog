package iost

import (
	"net"
	"net/http"
	"time"

	"github.com/EOSLaoMao/watchdog/internal/engine"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	StatusOK      string = "OK"
	StatusTimeout string = "Timeout"
	StatusFailure string = "Failure"
)

const nodeInfoURL = "http://iost.eoslaomao.com:30001/getNodeInfo"

func Listen() {
	engine.E().GET("/iost/status", func(c *gin.Context) {
		switch getNodeInfo() {
		case StatusOK:
			c.String(200, "IOST in good condition :), <i>%v</i>", time.Now().Format(time.RFC1123))
		case StatusTimeout:
			c.String(502, "request IOST node timeout, <i>%v</i>", time.Now().Format(time.RFC1123))
		case StatusFailure:
			c.String(502, "<b>IOST node maybe unavailable</b>, <i>%v</i>", time.Now().Format(time.RFC1123))
		}
	})
}

func getNodeInfo() string {
	res, err := http.Get(nodeInfoURL)
	if err != nil {
		logrus.Errorln("IOST Node Info Failure : ", err)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return StatusTimeout
		}
		return StatusFailure
	}

	if res.StatusCode != 200 {
		return StatusFailure
	}

	return StatusOK
}
