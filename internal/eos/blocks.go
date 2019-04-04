package eos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type BlockStatus string

const (
	BlockStatusPrepare BlockStatus = "Prepare"
	BlockStatusOK      BlockStatus = "OK"
	BlockStatusDown    BlockStatus = "Service maybe unavailable"
)

const blocksURL string = "https://api.eoslaomao.com/v1/chain/get_producers"

type BlockStat struct {
	UnpaidBlocks int
	Ranking      int
	Status       BlockStatus
}

type info struct {
	UnpaidBlocks int
	Ranking      int
}

var c = make(chan info)
var bs = &BlockStat{
	Status: BlockStatusPrepare,
}

func CheckUnpaidBlocks() {
	go func() {
		// 12*0.5*21
		ticker := time.NewTicker(BlockProduceTime * bpCount * time.Second)
		for range ticker.C {
			i, err := getUnpaidBlocks(bpName)
			if err != nil {
				fmt.Println("get unpaid blocks error: ", err.Error())
				bs.Status = BlockStatusDown
				continue
			}

			c <- *i
		}
	}()

	for {
		select {
		case i := <-c:
			switch {
			case i.UnpaidBlocks-bs.UnpaidBlocks >= 12:
				fallthrough
			case i.UnpaidBlocks-bs.UnpaidBlocks < 0:
				bs.Status = BlockStatusOK
				break
			default:
				bs.Status = BlockStatusDown
				break
			}

			bs.UnpaidBlocks = i.UnpaidBlocks
			bs.Ranking = i.Ranking
		}
	}
}

func getUnpaidBlocks(owner string) (*info, error) {
	cli := http.Client{
		Timeout: time.Duration(5 * time.Second),
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
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf(string(body))
	}

	type tableRows struct {
		Rows []struct {
			Owner        string `json:"owner"`
			UnpaidBlocks int    `json:"unpaid_blocks"`
		} `json:"rows"`
	}

	var rows tableRows
	_ = json.Unmarshal(body, &rows)

	for i, r := range rows.Rows {
		if r.Owner == owner {
			return &info{
				UnpaidBlocks: r.UnpaidBlocks,
				Ranking:      i + 1,
			}, nil
		}
	}

	return nil, fmt.Errorf("bp with name %s NOT FOUND", owner)
}
