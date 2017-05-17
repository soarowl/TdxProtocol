package test

import (
	. "github.com/onsi/ginkgo"

	"github.com/TdxProtocol/network"
	"fmt"
	"sort"
	"time"
)

var _ = Describe("BizApiGetSZStockCodes", func () {
	It("test", func() {
		fmt.Println("test GetSZStockCodes...")

		err, api := network.CreateBizApi(HOST_ONLY)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer api.Cleanup()

		start := time.Now().UnixNano()
		_, result := api.GetSHStockCodes()
		fmt.Println("got:", len(result), "time cost:", (time.Now().UnixNano() - start) / 1000000, "ms")
		sort.Strings(result)

		for _, c := range result {
			fmt.Println(c)
		}
	})
})

var _ = Describe("BizApiGetInfoEx", func () {
	It("test", func() {
		fmt.Println("test GetInfoEx...")
		err, api := network.CreateBizApi(HOST_ONLY)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer api.Cleanup()

		_, codes := api.GetSZStockCodes()

		start := time.Now().UnixNano()
		_, result := api.GetInfoEx(codes)
		fmt.Println("got:", len(result), "time cost:", (time.Now().UnixNano() - start) / 1000000, "ms")
		for k, l := range result {
			fmt.Println(k)
			for _, t := range l {
				fmt.Println(t)
			}
		}
	})
})

var _ = Describe("BizApiGetDayData", func () {
	It("test", func() {
		fmt.Println("test GetDayData...")
		err, api := network.CreateBizApi(HOST_ONLY)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer api.Cleanup()

		start := time.Now().UnixNano()
		_, result := api.GetLatestDayData("600000", 500)
		fmt.Println("got:", len(result), "time cost:", (time.Now().UnixNano() - start) / 1000000, "ms")
		for _, t := range result {
			fmt.Println(t)
		}
	})
})

var _ = Describe("BizApiMinuteDataPerf", func () {
	It("test", func() {
		fmt.Println("test BizApiMinuteDataPerf...")
		err, api := network.CreateBizApi(HOST_ONLY)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer api.Cleanup()

		_, codes := api.GetAStockCodes()

		//codes = codes[:10]

		doneChans := make([]chan int, 5)
		recordCh := make(chan map[string]interface{}, len(codes) + 1)

		count := (len(codes) + 4) / 5
		start := time.Now().UnixNano()
		for i := 0; i < 5; i++ {
			doneChans[i] = make(chan int)

			start := i * count
			end := (i + 1) * count
			if end > len(codes) {
				end = len(codes)
			}

			go func(codes []string, doneCh chan int) {
				for _, code := range codes {
					_, result := api.GetLatestMinuteData(code, 5)
					recordCh <- map[string]interface{}{"code": code, "record": result}
				}
				doneCh <- 1
			}(codes[start:end], doneChans[i])
		}

		for i := 0; i < 5; i++ {
			_ = <- doneChans[i]
			close(doneChans[i])
		}
		fmt.Println("time cost:", (time.Now().UnixNano() - start) / 1000000, "ms")

		recordCh <- map[string]interface{}{"code": ""}

		for {
			d := <- recordCh
			if d["code"] == "" {
				break
			}

			result, ok := d["record"].([]*network.Record)
			if ok {
				fmt.Println(d["code"], result[0])
			}
		}

		close(recordCh)
	})
})
