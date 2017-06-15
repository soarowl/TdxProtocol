package main

import (
	"github.com/stephenlyu/TdxProtocol/network"
	"fmt"
	"sort"
	"time"
	"encoding/json"
	"io/ioutil"
	"flag"
	"github.com/stephenlyu/TdxProtocol/util"
	"strings"
)


const (
	HOST = "125.39.80.98"
)


type record struct {
	Date string				`json:"date"`
	Open float32			`json:"open"`
	Close float32			`json:"close"`
	High float32			`json:"high"`
	Low float32				`json:"low"`
	Volume float32			`json:"volume"`
	Amount float32			`json:"amount"`
}


func fetchLatestMinuteData(host string, offset int, n int, codes []string) (error, map[string][]*record) {
	err, api := network.CreateBizApi(host)
	if err != nil {
		return err, nil
	}
	defer api.Cleanup()

	api.SetTimeOut(3 * 1000)

	if codes == nil {
		_, codes = api.GetAStockCodes()
	}
	sort.Strings(codes)
	nThread := len(codes) / 5
	if nThread == 0 {
		nThread = 1
	} else if (nThread > 10) {
		nThread = 10
	}

	doneChans := make([]chan int, nThread)
	recordCh := make(chan map[string]interface{}, len(codes) + 1)

	count := (len(codes) + 4) / nThread
	for i := 0; i < nThread; i++ {
		doneChans[i] = make(chan int)

		start := i * count
		end := (i + 1) * count
		if end > len(codes) {
			end = len(codes)
		}

		go func(codes []string, doneCh chan int) {
			for _, code := range codes {
				_, result := api.GetLatestMinuteData(code, offset, n)
				recordCh <- map[string]interface{}{"code": code, "record": result}
			}
			doneCh <- 1
		}(codes[start:end], doneChans[i])
	}

	for i := 0; i < nThread; i++ {
		_ = <- doneChans[i]
		close(doneChans[i])
	}

	recordCh <- map[string]interface{}{"code": ""}

	result := map[string][]*record{}

	for {
		d := <- recordCh
		code, _ := d["code"].(string)
		if code == "" {
			break
		}

		records, ok := d["record"].([]*network.Record)
		if ok {
			transRecords := make([]*record, len(records))
			for i, r := range records {
				tr := &record{
					Date: util.FormatMinuteDate(util.ToWindMinuteDate(r.Date)),
					Open: float32(r.Open) / 1000,
					Close: float32(r.Close) / 1000,
					High: float32(r.High) / 1000,
					Low: float32(r.Low) / 1000,
					Amount: r.Amount,
					Volume: r.Volume,
				}
				transRecords[i] = tr
			}

			result[code] = transRecords
		}
	}

	close(recordCh)

	return nil, result
}

func main() {
	host := flag.String("host", HOST, "服务器地址")
	stockCode := flag.String("stock-code", "", "股票代码，以都好分割")
	count := flag.Int("count", 5, "获取最近的K线数量")
	offset := flag.Int("offset", 0, "从倒数第一根K线开始获取")
	filePath := flag.String("output", "./minute-data.json", "文件名")
	flag.Parse()

	var stockCodes []string = nil
	if *stockCode != "" {
		stockCodes = strings.Split(*stockCode, ",")
	}

	start := time.Now().UnixNano()

	var err error
	var result map[string][]*record
	for {
		err, result = fetchLatestMinuteData(*host, *offset, *count, stockCodes)
		if err != nil {
			fmt.Println("try get minute data error", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}
	fmt.Println("time cost:", (time.Now().UnixNano() - start) / 1000000, "ms")
	fmt.Println("total: ", len(result))

	bytes, _ := json.MarshalIndent(result, "", "  ")
	err = ioutil.WriteFile(*filePath, bytes, 0666)
	if err != nil {
		panic(err)
	}
}
