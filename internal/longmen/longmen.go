package longmen

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/EOSLaoMao/watchdog/internal/engine"
)

const (
	StatusOK      string = "OK"
	StatusTimeout string = "Timeout"
	StatusFailure string = "Failure"
)

const ServerInfoURL = "https://longmen.fund/api/price"

func Listen() {
	engine.E().GET("/longmen/status", func(c *gin.Context) {
		switch getServerInfo() {
		case StatusOK:
			c.String(200, "OK")
		case StatusTimeout:
			c.String(504, "request Longmen server timeout")
		case StatusFailure:
			c.String(502, "<b>Longmen server maybe down</b>")
		}
	})
}

func getServerInfo() string {
	cli := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	req, _ := http.NewRequest("GET", ServerInfoURL, nil)

	res, err := cli.Do(req)
	if err != nil {
		logrus.Errorln("Longmen Failure : ", err)
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
