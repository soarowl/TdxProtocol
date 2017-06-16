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
	"path/filepath"
)


const (
	HOST = "125.39.80.98"
	MAX_THREAD = 20
	DAY = time.Hour * 24
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
	} else if (nThread > MAX_THREAD) {
		nThread = MAX_THREAD
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

func saveData(filePath string, result map[string][]*record) {
	bytes, _ := json.MarshalIndent(result, "", "  ")
	err := ioutil.WriteFile(filePath, bytes, 0666)
	if err != nil {
		panic(err)
	}
}

func tryFetchData(host string, offset int, count int, stockCodes []string) map[string][]*record {
	var err error
	var result map[string][]*record
	for {
		err, result = fetchLatestMinuteData(host, offset, count, stockCodes)
		if err != nil {
			fmt.Println("try get minute data error", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	return result
}

func runOnce(host string, offset int, count int, stockCodes []string, filePath string) {
	start := time.Now().UnixNano()

	result := tryFetchData(host, offset, count, stockCodes)

	fmt.Println("time cost:", (time.Now().UnixNano() - start) / 1000000, "ms")
	fmt.Println("total: ", len(result))

	saveData(filePath, result)
}

func runOnceForDaemon(date string, host string, offset int, stockCodes []string, filePath string) {
	start := time.Now().UnixNano()

	result := tryFetchData(host, offset, 1, stockCodes)

	if offset == 1 {
		codes := []string{}
		for code, record := range result {
			if len(record) == 0 {
				continue
			}

			if record[0].Date == date {
				continue
			}

			codes = append(codes, code)
		}
		if len(codes) > 0 {
			fmt.Println(codes)
		}

		pResult := tryFetchData(host, 0, 2, codes)
		for code, record := range pResult {
			if len(record) == 0 {
				continue
			}

			if record[0].Date == date {
				result[code] = record[0:1]
			} else if record[1].Date == date {
				result[code] = record[1:]
			}
		}
	}
	fmt.Println("time cost:", (time.Now().UnixNano() - start) / 1000000, "ms", " total: ", len(result))

	saveData(filePath, result)
}

func main() {
	host := flag.String("host", HOST, "服务器地址")
	stockCode := flag.String("stock-code", "", "股票代码，以都好分割")
	count := flag.Int("count", 5, "获取最近的K线数量")
	offset := flag.Int("offset", 0, "从倒数第一根K线开始获取")
	filePath := flag.String("output", "./minute-data.json", "文件名")
	daemon := flag.Bool("daemon", false, "是否一直运行")

	flag.Parse()

	var stockCodes []string = nil
	if *stockCode != "" {
		stockCodes = strings.Split(*stockCode, ",")
	}

	if *daemon {
		var lastDate string
		var prevDate string
		var offset int

		dirName := filepath.Dir(*filePath)
		fileName := filepath.Base(*filePath)
		extName := filepath.Ext(fileName)
		mainName := fileName[0:len(fileName) - len(extName)]

		indexCode := "999999"
		for {
			minute := util.GetTimeString()
			if minute > "15:01:00" {
				now := time.Now()
				_, zoneOffset := now.Zone()
				timestamp := now.UnixNano() + int64(zoneOffset) * int64(time.Second)

				sleepDuration := DAY - time.Duration(timestamp) % DAY + 9 * time.Hour + 25 * time.Minute
				nextTradeDayStart := now.Add(sleepDuration)
				fmt.Printf("Sleep until %s.\n", util.FormatLongDate(nextTradeDayStart))
				time.Sleep(sleepDuration)
			} else if minute < "09:25:00" {
				now := time.Now()
				_, zoneOffset := now.Zone()
				timestamp := now.UnixNano() + int64(zoneOffset) * int64(time.Second)

				sleepDuration := 9 * time.Hour + 25 * time.Minute - time.Duration(timestamp) % DAY
				nextTradeDayStart := now.Add(sleepDuration)
				fmt.Printf("Sleep until %s.\n", util.FormatLongDate(nextTradeDayStart))
				if sleepDuration > 0 {
					time.Sleep(sleepDuration)
				}
			}

			today := util.GetTodayString()

			loopStart := time.Now().UnixNano()
			for {
				err, result := fetchLatestMinuteData(*host, 0, 1, []string{indexCode})
				if err != nil {
					fmt.Println("loop fail, error:", err)
					time.Sleep(500 * time.Millisecond)
					continue
				}
				if _, ok := result[indexCode]; !ok {
					continue
				}
				if len(result[indexCode]) == 0 {
					continue
				}
				r := result[indexCode][0]

				if r.Date[0:8] != today {
					time.Sleep(1 * time.Second)
					continue
				}

				if lastDate != "" && r.Date != lastDate {
					lastDate = r.Date
					offset = 1
					time.Sleep(10 * time.Second)		// Sleep 10 seconds to wait all stock has the predate data.
					break
				}

				if r.Date[9:] == "14:59:00" && time.Now().UnixNano() - loopStart > int64(time.Minute + 2 * time.Second) {
					prevDate = lastDate
					lastDate = r.Date
					offset = 0
					break
				}

				lastDate = r.Date
				if prevDate == "" {
					prevDate = lastDate
				}
			}

			outputFile := filepath.Join(dirName, fmt.Sprintf("%s-%s%s", mainName, prevDate[9:], extName))
			fmt.Printf("Fetching data at %s...\n", lastDate)
			runOnceForDaemon(prevDate, *host, offset, stockCodes, outputFile)

			prevDate = lastDate
		}
	} else {
		runOnce(*host, *offset, *count, stockCodes, *filePath)
	}
}
