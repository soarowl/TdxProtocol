package main

import (
	"flag"
	"github.com/stephenlyu/TdxProtocol/network"
	"fmt"
	"os"
)

const (
	HOST = "125.39.80.98"
)

func main() {
	output := flag.String("output", "temp", "Directory save file to")
	flag.Parse()

	err, api := network.CreateBizApi(HOST)
	if err != nil {
		fmt.Errorf("[ERROR] Connect server fail, error: %+v", err)
		os.Exit(1)
	}
	defer api.Cleanup()

	for _, fileName := range flag.Args() {
		err = api.DownloadFile(fileName, *output)
		if err != nil {
			fmt.Errorf("[ERROR] Download %s fail, error: %+v", err)
		}
	}
}
