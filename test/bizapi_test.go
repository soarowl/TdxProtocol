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
		_, result := api.GetSZStockCodes()
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
