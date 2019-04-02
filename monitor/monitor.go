package monitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Monitor() {
	ticker := time.NewTicker(6 * time.Second)
	for range ticker.C {
		res, err := http.Get("http://127.0.0.1:8080/block/status")
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
	}
}
