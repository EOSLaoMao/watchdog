package mixin

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/EOSLaoMao/watchdog/internal/engine"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var mixinPaths = map[string]string{
	"node0": "http://mixin-node0.eoslaomao.com:1443",
	"node1": "http://mixin-node1.eoslaomao.com:1443",
}

const (
	StatusOK      string = "OK"
	StatusTimeout string = "Timeout"
	StatusFailure string = "Failure"
)

func Listen() {
	for n, u := range mixinPaths {
		serve(n, u)
	}
}

func serve(name string, url string) {
	engine.E().GET(fmt.Sprintf("/mixin/%s/status", name), func(c *gin.Context) {
		switch getMixinInfo(url) {
		case StatusOK:
			c.String(200, "<b>%s</b> in good condition :), <i>%v</i>", name, time.Now().Format(time.RFC1123))
		case StatusTimeout:
			c.String(502, "request timeout for <b>%s</b>, <i>%v</i>", name, time.Now().Format(time.RFC1123))
		case StatusFailure:
			c.String(502, "<b>%s maybe unavailable</b>, <i>%v</i>", name, time.Now().Format(time.RFC1123))
		}
	})
}

func getMixinInfo(url string) string {
	json := []byte(`{"method": "getinfo"}`)
	body := bytes.NewBuffer(json)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, body)

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorln("Mixin Failure : ", err)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return StatusTimeout
		}
		return StatusFailure
	}

	if resp.StatusCode != 200 {
		return StatusFailure
	}

	return StatusOK
}
