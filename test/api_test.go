package test

import (
	. "github.com/onsi/ginkgo"

	"github.com/stephenlyu/TdxProtocol/network"
	"fmt"
	"sort"
)

var _ = Describe("GetStockList", func () {
	It("test", func() {
		err, api := network.CreateAPI(HOST)
		if err != nil {
			chk(err)
			return
		}
		defer api.Cleanup()

		err, total, result := api.GetStockList(network.BLOCK_SH_A, 0, 10)
		chk(err)
		fmt.Println("total:", total, " got:", len(result))
		codes := []string{}
		for k, _ := range result {
			codes = append(codes, k)
		}
		sort.Strings(codes)

		for _, c := range codes {
			fmt.Println(c)
		}
	})
})


var _ = Describe("GetInfoEx", func () {
	It("test", func() {
		err, api := network.CreateAPI(HOST)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer api.Cleanup()

		_, result := api.GetInfoEx([]string{"600000", "000001"})
		for k, l := range result {
			fmt.Println(k)
			for _, t := range l {
				fmt.Println(t)
			}
		}
	})
})

var _ = Describe("GetMinuteData", func () {
	It("test", func() {
		err, api := network.CreateAPI(HOST)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer api.Cleanup()

		_, result := api.GetMinuteData("600000", 0, 10)
		for _, t := range result {
			fmt.Println(t)
		}
	})
})
