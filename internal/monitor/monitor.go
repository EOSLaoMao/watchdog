package monitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/EOSLaoMao/watchdog/internal/eos"
	"github.com/EOSLaoMao/watchdog/internal/message"
)

var monitorList = map[string]string{
	"EOS":          eos.ListenBlockPath,
	"Mixin node 0": "/mixin/node0/status",
	"Mixin node 1": "/mixin/node1/status",
	"IOST":         "/iost/status",
}

func StartMonitor() {
	ticker := time.NewTicker(3 * time.Minute)
	for range ticker.C {

		for k, v := range monitorList {
			host := os.Getenv("MONITOR_SERVER")
			res, err := http.Get(fmt.Sprintf("http://%s:8080%s", host, v))
			if err != nil {
				message.SendToTelegram(
					fmt.Sprintf("Monitor for %s maybe unavailable: %s", k, err.Error()),
				)
				break
			}
			defer res.Body.Close()

			body, _ := ioutil.ReadAll(res.Body)
			msg := fmt.Sprintf("%s: %s", k, string(body))

			switch res.StatusCode {
			case 200:
				message.SendToTelegram(msg)
			case 502:
				message.SendToTelegram(msg)
			}
		}
	}
}
