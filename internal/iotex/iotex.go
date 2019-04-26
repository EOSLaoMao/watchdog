package iotex

import (
	"io/ioutil"
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

const nodeInfoURL = "http://iotex-http.eoslaomao.com/health"

func Listen() {
	engine.E().GET("/iotex/status", func(c *gin.Context) {
		switch getNodeInfo() {
		case StatusOK:
			c.String(200, "IoTex in good condition :), <i>%v</i>", time.Now().Format(time.RFC1123))
		case StatusTimeout:
			c.String(502, "request IoTex node timeout, <i>%v</i>", time.Now().Format(time.RFC1123))
		case StatusFailure:
			c.String(502, "<b>IoTex node maybe unavailable</b>, <i>%v</i>", time.Now().Format(time.RFC1123))
		}
	})
}

func getNodeInfo() string {
	cli := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	req, _ := http.NewRequest("GET", nodeInfoURL, nil)

	res, err := cli.Do(req)
	if err != nil {
		logrus.Errorln("IoTex Node Info Failure : ", err)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return StatusTimeout
		}
		return StatusFailure
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if string(body) != "OK" {
		return StatusFailure
	}

	return StatusOK
}