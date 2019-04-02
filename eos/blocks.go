package eos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const blocksURL string = "https://api.eoslaomao.com/v1/chain/get_table_rows"

type blocks struct {
	Blocks []uint
	sync.RWMutex
}

type blockStatus struct {
	Err          error
	UnpaidBlocks uint
	Time         time.Time
}

var b = &blocks{}
var s = &blockStatus{}

func (bs *blockStatus) statusError(err error) {
	bs.Err = err
	bs.Time = time.Now()
}

func logBlocks() {
	// 12*0.5*21
	ticker := time.NewTicker(BlockProduceTime * time.Second)
	for range ticker.C {
		blocks, err := unpaidBlocks(bpName)
		if err != nil {
			s.Time = time.Now()
			continue
		}
		s.UnpaidBlocks = blocks
		s.Time = time.Now()

		b.Lock()

		b.Blocks = append(b.Blocks, blocks)
		b.Unlock()
	}
}

func checkBlocks() {
	ticker := time.NewTicker(BlockProduceTime * bpCount * time.Second)
	for range ticker.C {
		if len(b.Blocks) == 0 {
			continue
		}

		b.RLock()
		minus := b.Blocks[len(b.Blocks)-1] - b.Blocks[0]
		if minus >= 12 || minus < 0 {
			s.Err = nil
		} else {
			s.Err = fmt.Errorf("service maybe unavailable")
		}
		s.Time = time.Now()
		b.RUnlock()

		b.Lock()
		b.Blocks = b.Blocks[:0]
		b.Unlock()
	}
}

func unpaidBlocks(owner string) (uint, error) {
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
