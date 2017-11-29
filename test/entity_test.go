package test

import (
	. "github.com/onsi/ginkgo"

	"github.com/stephenlyu/TdxProtocol/network"
	"fmt"
	"bytes"
	"net"
	"encoding/hex"
	"sort"
	"time"
	"baiwenbao.com/arbitrage/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	HOST_ONLY = "125.39.80.98"
	HOST = "125.39.80.98:7709"
)

func chk(err error) {
	if err == nil {
		return
	}

	fmt.Println("error:", err)
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
		_, result := parser.Parse()
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
		_, result := parser.Parse()
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
	req := network.NewInstantTransReq(1, "600000", 0, 100)
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
		_, result := parser.Parse()
		//fmt.Println(hex.EncodeToString(parser.Data))

		fmt.Println("record count: ", len(result))
		for _, t := range result {
			fmt.Println(t)
		}
	})
})

func BuildHisTransBuffer() (*bytes.Buffer, *network.HisTransReq) {
	req := network.NewHisTransReq(1, 20170414, "600000", 0, 100)
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
		_, result := parser.Parse()

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
		_, result := parser.Parse()
		fmt.Println(hex.EncodeToString(parser.Data))

		fmt.Println("record count: ", len(result))
		for _, t := range result {
			fmt.Println(t)
		}
	})
})

func BuildGetFileLenBuffer(fileName string) (*bytes.Buffer, *network.GetFileLenReq) {
	req := network.NewGetFileLenReq(1, fileName)
	buf := new(bytes.Buffer)
	req.Write(buf)
	return buf, req
}

var _ = Describe("TestGetFileLenReq", func() {
	It("test", func() {
		fmt.Println("TestGetFileLenReq...")
		buf, req := BuildGetFileLenBuffer("zhb.zip")

		start := time.Now().UnixNano()
		conn, err := net.Dial("tcp", HOST)
		chk(err)
		defer conn.Close()

		_, err = conn.Write(buf.Bytes())
		chk(err)

		err, buffer := network.ReadResp(conn)
		chk(err)
		fmt.Println("time cost: ", time.Now().UnixNano() - start)

		parser := network.NewGetFileLenParser(req, buffer)
		_, length := parser.Parse()
		fmt.Println(hex.EncodeToString(parser.Data))

		fmt.Println("file length: ", length)
	})
})

func BuildGetFileDataBuffer(fileName string, offset uint32, length uint32) (*bytes.Buffer, *network.GetFileDataReq) {
	req := network.NewGetFileDataReq(1, fileName, offset, length)
	buf := new(bytes.Buffer)
	req.Write(buf)
	return buf, req
}

var _ = Describe("TestGetFileDataReq", func() {
	var getFileLen = func (fileName string) uint32 {
		buf, req := BuildGetFileLenBuffer(fileName)

		start := time.Now().UnixNano()
		conn, err := net.Dial("tcp", HOST)
		chk(err)
		defer conn.Close()

		_, err = conn.Write(buf.Bytes())
		chk(err)

		err, buffer := network.ReadResp(conn)
		chk(err)
		fmt.Println("time cost: ", time.Now().UnixNano() - start)

		parser := network.NewGetFileLenParser(req, buffer)
		_, length := parser.Parse()
		return length
	}

	It("test", func() {
		fmt.Println("TestGetFileDataReq...")
		fileName := "bi/bigdata.zip"

		length := getFileLen(fileName)

		start := time.Now().UnixNano()
		conn, err := net.Dial("tcp", HOST)
		chk(err)
		defer conn.Close()

		fileData := make([]byte, length)

		var offset uint32 = 0
		var count uint32 = 30000

		for offset < length {
			buf, req := BuildGetFileDataBuffer(fileName, offset, count)
			_, err = conn.Write(buf.Bytes())
			chk(err)

			err, buffer := network.ReadResp(conn)
			chk(err)
			fmt.Println("time cost: ", time.Now().UnixNano() - start)

			parser := network.NewGetFileDataParser(req, buffer)
			_, packetLength, data := parser.Parse()
			util.Assert(packetLength == uint32(len(data)), "")

			copy(fileData[offset:offset + packetLength], data[:])

			offset += count
		}

		filePath := filepath.Join("temp", fileName)
		os.MkdirAll(filepath.Dir(filePath), 0777)
		ioutil.WriteFile(filePath, fileData, 0666)
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

var _ = Describe("TestCmd06b9", func () {
	It("test", func() {
		reqHex := "0c5e186a00016e006e00b906b01e0400307500007a68622e7a6970000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
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

		fmt.Printf("data len: %d\n", len(parser.Data))
		fmt.Println(hex.EncodeToString(parser.Data))
	})
})
