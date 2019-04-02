package monitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/EOSLaoMao/watchdog/eos"
	"github.com/EOSLaoMao/watchdog/message"
)

func Monitor() {
	var t int

	ticker := time.NewTicker(eos.BlockProduceTime * time.Second)
	for range ticker.C {
		t++

		res, err := http.Get("http://127.0.0.1:8080/block/status")
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		if res.StatusCode == 500 {
			message.SendToTelegram(string(body))
			continue
		}

		if t == 60 {
			message.SendToTelegram(string(body))
			t = 0
		}
	}
}
