package monitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/EOSLaoMao/watchdog/internal/eos"
	"github.com/EOSLaoMao/watchdog/internal/message"
)

const alertTimes = 3

var monitorList = map[string]string{
	"EOS":          eos.ListenBlockPath,
	"Mixin node 0": "/mixin/node0/status",
	"Mixin node 1": "/mixin/node1/status",
	"IoTex":        "/iotex/status",
	"Longmen":      "/longmen/status",
}

var timeoutCache = make(map[string]int)

type monitorMsg struct {
	Symbol string
	Code   int
	Msg    string
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
				Symbol: k,
				Code:   res.StatusCode,
				Msg:    fmt.Sprintf("<b>%s</b>: %s", k, string(body)),
			}

			msgs = append(msgs, msg)
		}

		sort.Slice(msgs, func(i, j int) bool {
			return msgs[i].Symbol <= msgs[j].Symbol
		})

		ok := true
		var result []string
		for _, m := range msgs {
			result = append(result, m.Msg)

			switch m.Code {
			case 200:
				timeoutCache[m.Symbol] = 0
			case 504:
				timeoutCache[m.Symbol]++
			case 502:
				//ok = false
				timeoutCache[m.Symbol]++
			}
		}

		for _, v := range timeoutCache {
			if v >= alertTimes {
				ok = false
			}
		}

		if !ok {
			message.MakeVoiceCall()
		}

		r := fmt.Sprintf(
			"%s\n\n<i>%v</i>",
			strings.Join(result, "\n\n"),
			time.Now().Format(time.RFC1123),
		)
		message.SendToTelegram(url.QueryEscape(r))
	}
}
