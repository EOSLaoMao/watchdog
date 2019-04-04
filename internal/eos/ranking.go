package eos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const rankingURL = "https://api.eoslaomao.com/v1/chain/get_producers"

var ranking int

func CheckBPRanking() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		r, err := getBPRanking()
		if err != nil {
			fmt.Println(err)
			continue
		}
		ranking = r
	}
}

func getBPRanking() (int, error) {
	cli := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	params := map[string]interface{}{
		"lower_bound": "",
		"limit":       100,
		"json":        true,
	}
	b, _ := json.Marshal(params)

	req, _ := http.NewRequest("POST", blocksURL, bytes.NewBuffer(b))

	res, err := cli.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return 0, fmt.Errorf(string(body))
	}

	var result []map[string]interface{}

	_ = json.Unmarshal(body, &result)

	for i, r := range result {
		if r["owner"] == bpName {
			return i, nil
		}
	}

	return 0, fmt.Errorf("bp with name %s NOT FOUND", bpName)
}
