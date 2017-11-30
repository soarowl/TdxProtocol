package test

import (
	. "github.com/onsi/ginkgo"

	"github.com/stephenlyu/TdxProtocol/network"
	"fmt"
)


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
