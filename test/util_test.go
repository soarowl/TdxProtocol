package test

import (
. "github.com/onsi/ginkgo"

	"fmt"
	"github.com/TdxProtocol/util"
)

var _ = Describe("FormatMinuteDate", func() {
	It("test", func (){
		fmt.Println(util.FormatMinuteDate(58747397))
	})
})