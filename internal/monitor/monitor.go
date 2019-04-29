package monitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/EOSLaoMao/watchdog/internal/eos"
	"github.com/EOSLaoMao/watchdog/internal/message"
)

var monitorList = map[string]string{
	"EOS":          eos.ListenBlockPath,
	"Mixin node 0": "/mixin/node0/status",
	"Mixin node 1": "/mixin/node1/status",
	"IoTex":        "/iotex/status",
}

type monitorMsg struct {
	Code int
	Msg  string
}

func StartMonitor() {
	ticker := time.NewTicker(3 * time.Minute)
	for range ticker.C {

		var msgs []*monitorMsg

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

			msg := &monitorMsg{
				Code: res.StatusCode,
				Msg:  fmt.Sprintf("<i>%s</i>: %s", k, string(body)),
			}

			msgs = append(msgs, msg)
		}

		var result []string
		for _, m := range msgs {
			result = append(result, m.Msg)
		}

		m := fmt.Sprintf("%s\n <i>%v</i>", strings.Join(result, "\n"), time.Now().Format(time.RFC1123))
		message.SendToTelegram(m)
	}
}
