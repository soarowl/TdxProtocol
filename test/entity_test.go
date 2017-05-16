package test

import (
	. "github.com/onsi/ginkgo"

	"github.com/TdxProtocol/network"
	"fmt"
	"bytes"
	"net"
	"encoding/hex"
	"sort"
	"time"
)

const (
	HOST_ONLY = "125.39.80.98"
	HOST = "125.39.80.98:7709"
)

func chk(err error) {
	if err == nil {
		return
	}

	fmt.Println(err)
	panic(err)
}

func BuildStockListBuffer() (*bytes.Buffer, *network.StockListReq) {
	req := network.NewStockListReq(1, 0, 80, 80)
	buf := new(bytes.Buffer)
	req.Write(buf)
	return buf, req
}

var _ = Describe("TestSockListReq", func() {
	It("test", func () {
		fmt.Println("TestStockListReq...")
		buf, req := BuildStockListBuffer()

		conn, err := net.Dial("tcp", HOST)
		chk(err)
		defer conn.Close()

		fmt.Println(hex.EncodeToString(buf.Bytes()))
		_, err = conn.Write(buf.Bytes())
		chk(err)

		err, buffer := network.ReadResp(conn)
		chk(err)

		parser := network.NewStockListParser(req, buffer)
		result := parser.Parse()
		fmt.Println(hex.EncodeToString(parser.Data))

		fmt.Println("total:", parser.Total, " got:", len(result))
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

func BuildInfoExBuffer() (*bytes.Buffer, *network.InfoExReq) {
	req := network.NewInfoExReq(1)
	req.AddCode("000099")
	fmt.Println(req)
	buf := new(bytes.Buffer)
	req.Write(buf)
	fmt.Println(buf.Bytes())
	return buf, req
}

var _ = Describe("TestInfoExReq", func() {
	It("test", func () {
		fmt.Println("TestInfoExReq...")
		buf, req := BuildInfoExBuffer()

		conn, err := net.Dial("tcp", HOST)
		chk(err)
		defer conn.Close()

		_, err = conn.Write(buf.Bytes())
		chk(err)

		err, buffer := network.ReadResp(conn)
		chk(err)

		parser := network.NewInfoExParser(req, buffer)
		result := parser.Parse()
		fmt.Println(hex.EncodeToString(parser.Data))

		for k, l := range result {
			fmt.Println(k)
			for _, t := range l {
				fmt.Println(t)
			}
		}
	})
})

func BuildInstantTransBuffer() (*bytes.Buffer, *network.InstantTransReq){
	req := network.NewInstantTransReq(1, "600000", 4000, 6000)
	buf := new(bytes.Buffer)
	req.Write(buf)
	return buf, req
}

var _ = Describe("TestInstantTransReq", func() {
	It("test", func() {
		fmt.Println("TestInstantTransReq...")
		buf, req := BuildInstantTransBuffer()

		conn, err := net.Dial("tcp", HOST)
		chk(err)
		defer conn.Close()

		_, err = conn.Write(buf.Bytes())
		chk(err)

		err, buffer := network.ReadResp(conn)
		chk(err)

		parser := network.NewInstantTransParser(req, buffer)
		result := parser.Parse()
		//fmt.Println(hex.EncodeToString(parser.Data))

		fmt.Println("record count: ", len(result))
		for _, t := range result {
			fmt.Println(t)
		}
	})
})

func BuildHisTransBuffer() (*bytes.Buffer, *network.HisTransReq) {
	req := network.NewHisTransReq(1, 20170414, "600000", 2000, 1)
	buf := new(bytes.Buffer)
	req.Write(buf)
	return buf, req
}

var _ = Describe("TestHisTransReq", func() {
	It("test", func() {
		fmt.Println("TestHisTransReq...")
		buf, req := BuildHisTransBuffer()

		start := time.Now().UnixNano()
		conn, err := net.Dial("tcp", HOST)
		chk(err)
		defer conn.Close()

		_, err = conn.Write(buf.Bytes())
		chk(err)

		err, buffer := network.ReadResp(conn)
		chk(err)
		fmt.Println("time cost: ", time.Now().UnixNano() - start)

		parser := network.NewHisTransParser(req, buffer)
		result := parser.Parse()

		fmt.Println("record count: ", len(result))
		for _, t := range result {
			fmt.Println(t)
		}
	})
})

func BuildPeriodDataBuffer() (*bytes.Buffer, *network.PeriodDataReq) {
	req := network.NewPeriodDataReq(1, "600000", network.PERIOD_DAY, 0, 0x118)
	buf := new(bytes.Buffer)
	req.Write(buf)
	return buf, req
}

var _ = Describe("TestPeriodDataReq", func() {
	It("test", func() {
		fmt.Println("TestPeriodDataReq...")
		buf, req := BuildPeriodDataBuffer()

		start := time.Now().UnixNano()
		conn, err := net.Dial("tcp", HOST)
		chk(err)
		defer conn.Close()

		_, err = conn.Write(buf.Bytes())
		chk(err)

		err, buffer := network.ReadResp(conn)
		chk(err)
		fmt.Println("time cost: ", time.Now().UnixNano() - start)

		parser := network.NewPeriodDataParser(req, buffer)
		result := parser.Parse()
		fmt.Println(hex.EncodeToString(parser.Data))

		fmt.Println("record count: ", len(result))
		for _, t := range result {
			fmt.Println(t)
		}
	})
})

var _ = Describe("TestReqData", func () {
	It("test", func() {
		reqHex := "0c57086401011c001c002d050000333030313732040001000100180100000000000000000000"
		reqData, _ := hex.DecodeString(reqHex)

		conn, err := net.Dial("tcp", HOST)
		chk(err)
		defer conn.Close()

		_, err = conn.Write(reqData)
		chk(err)

		err, buffer := network.ReadResp(conn)
		chk(err)

		parser := network.NewRespParser(buffer)
		parser.Parse()

		fmt.Println(hex.EncodeToString(parser.Data))
	})
})
