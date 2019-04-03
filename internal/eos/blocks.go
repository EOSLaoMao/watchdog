package eos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Status string

const (
	StatusPrepare Status = "Prepare"
	StatusOK      Status = "OK"
	StatusDown    Status = "Service maybe unavailable"
)

const blocksURL string = "https://api.eoslaomao.com/v1/chain/get_table_rows"

type BlockStatus struct {
	UnpaidBlocks uint
	Status       Status
}

var c = make(chan uint)
var bs = &BlockStatus{
	Status: StatusPrepare,
}

func CheckUnpaidBlocks() {
	go func() {
		// 12*0.5*21
		ticker := time.NewTicker(BlockProduceTime * bpCount * time.Second)
		for range ticker.C {
			blocks, err := getUnpaidBlocks(bpName)
			if err != nil {
				bs.Status = StatusDown
				continue
			}

			c <- blocks
		}
	}()

	for {
		select {
		case blocks := <-c:
			switch {
			case blocks-bs.UnpaidBlocks >= 12:
				fallthrough
			case blocks-bs.UnpaidBlocks < 0:
				bs.Status = StatusOK
				break
			default:
				bs.Status = StatusDown
				break
			}

			bs.UnpaidBlocks = blocks
		}
	}
}

func getUnpaidBlocks(owner string) (uint, error) {
	cli := http.Client{
		Timeout: time.Duration(1 * time.Second),
	}

	params := map[string]interface{}{
		"json":        true,
		"code":        "eosio",
		"scope":       "eosio",
		"table":       "producers",
		"table_key":   "",
		"lower_bound": "",
		"upper_bound": "",
		"limit":       1000,
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

	type tableRows struct {
		Rows []struct {
			Owner        string `json:"owner"`
			UnpaidBlocks uint   `json:"unpaid_blocks"`
		} `json:"rows"`
	}

	var rows tableRows
	_ = json.Unmarshal(body, &rows)

	for _, r := range rows.Rows {
		if r.Owner == owner {
			return r.UnpaidBlocks, nil
		}
	}

	return 0, fmt.Errorf("bp with name %s NOT FOUND", owner)
}
