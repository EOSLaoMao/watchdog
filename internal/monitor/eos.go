package monitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/EOSLaoMao/watchdog/internal/eos"
	"github.com/EOSLaoMao/watchdog/internal/message"
)

var monitorList = map[string]string{
	"EOS": eos.ListenPath,
}

func StartMonitor() {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {

		for k, v := range monitorList {

			res, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080%s", v))
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
