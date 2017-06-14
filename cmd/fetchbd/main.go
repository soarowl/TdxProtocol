package main

import (
	"github.com/stephenlyu/TdxProtocol/network"
	"fmt"
	"sort"
	"time"
	"encoding/json"
	"io/ioutil"
	"flag"
)


const (
	HOST = "125.39.80.98"
)

func getStockCodes(api *network.BizApi) (err error, stockCodes []string) {

	err, codes := api.GetSHStockCodes()
	if err != nil {
		return
	}

	stockCodes = append(stockCodes, codes...)

	err, codes = api.GetSZStockCodes()
	if err != nil {
		return
	}

	stockCodes = append(stockCodes, codes...)
	return
}

func getInfoEx(api *network.BizApi) (error, map[string][]*network.InfoExItem) {
	err, codes := getStockCodes(api)
	sort.Strings(codes)
	if err != nil {
		return err, nil
	}

	return api.GetInfoEx(codes)
}

func tryGetInfoEx(host string) (error, map[string][]*network.InfoExItem) {
	err, api := network.CreateBizApi(host)
	if err != nil {
		panic(err)
	}
	defer api.Cleanup()

	return getInfoEx(api)
}

func main() {
	host := flag.String("host", HOST, "服务器地址")
	filePath := flag.String("output", "./info_ex.json", "文件名")
	flag.Parse()

	var err error
	var result map[string][]*network.InfoExItem
	for {
		err, result = tryGetInfoEx(*host)
		if err != nil {
			fmt.Println("try get info ex error", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	finalResult := map[string]interface{}{}

	for code, items := range result {
		var market string
		switch code[:2] {
		case "60":
			market = "sh"
		case "00":
			fallthrough
		case "30":
			market = "sz"
		default:
			continue
		}

		if _, ok := finalResult[market]; !ok {
			finalResult[market] = map[string]interface{}{}
		}

		marketValues, _ := finalResult[market]

		infoEx := marketValues.(map[string]interface{})

		infoEx[code] = map[string]interface{}{
			"info": map[string]string{},
			"ex": items,
		}
	}

	bytes, _ := json.MarshalIndent(finalResult, "", "  ")
	err = ioutil.WriteFile(*filePath, bytes, 0666)
	if err != nil {
		panic(err)
	}
}
