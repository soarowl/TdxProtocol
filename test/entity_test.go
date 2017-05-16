package test

import (
	"github.com/TdxProtocol/entity"
	"fmt"
	"bytes"
	"net"
	"encoding/hex"
	"testing"
	"sort"
	"time"
)

const (
	HOST = "125.39.80.98:7709"
)

func chk(err error) {
	if err == nil {
		return
	}

	fmt.Println(err)
	panic(err)
}

func BuildStockListBuffer() (*bytes.Buffer, *entity.StockListReq) {
	req := entity.NewStockListReq(1, 0, 80, 80)
	buf := new(bytes.Buffer)
	req.Write(buf)
	return buf, req
}

func _TestStockListReq(t *testing.T) {
	fmt.Println("TestStockListReq...")
	buf, req := BuildStockListBuffer()

	conn, err := net.Dial("tcp", HOST)
	chk(err)

	fmt.Println(hex.EncodeToString(buf.Bytes()))
	_, err = conn.Write(buf.Bytes())
	chk(err)

	err, buffer := entity.ReadResp(conn)
	chk(err)

	parser := entity.NewStockListParser(req, buffer)
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
}

func BuildInfoExBuffer() (*bytes.Buffer, *entity.InfoExReq) {
	req := entity.NewInfoExReq(1)
	req.AddCode("000099")
	fmt.Println(req)
	buf := new(bytes.Buffer)
	req.Write(buf)
	fmt.Println(buf.Bytes())
	return buf, req
}

func _TestInfoExReq(t *testing.T) {
	fmt.Println("TestInfoExReq...")
	buf, req := BuildInfoExBuffer()

	conn, err := net.Dial("tcp", HOST)
	chk(err)

	_, err = conn.Write(buf.Bytes())
	chk(err)

	err, buffer := entity.ReadResp(conn)
	chk(err)

	parser := entity.NewInfoExParser(req, buffer)
	result := parser.Parse()
	fmt.Println(hex.EncodeToString(parser.Data))

	for k, l := range result {
		fmt.Println(k)
		for _, t := range l {
			fmt.Println(t)
		}
	}
}

func BuildInstantTransBuffer() (*bytes.Buffer, *entity.InstantTransReq){
	req := entity.NewInstantTransReq(1, "600000", 4000, 6000)
	buf := new(bytes.Buffer)
	req.Write(buf)
	return buf, req
}

func _TestInstantTransReq(t *testing.T) {
	fmt.Println("TestInstantTransReq...")
	buf, req := BuildInstantTransBuffer()

	conn, err := net.Dial("tcp", HOST)
	chk(err)

	_, err = conn.Write(buf.Bytes())
	chk(err)

	err, buffer := entity.ReadResp(conn)
	chk(err)

	parser := entity.NewInstantTransParser(req, buffer)
	result := parser.Parse()
	//fmt.Println(hex.EncodeToString(parser.Data))

	fmt.Println("record count: ", len(result))
	for _, t := range result {
		fmt.Println(t)
	}
}

func BuildHisTransBuffer() (*bytes.Buffer, *entity.HisTransReq) {
	req := entity.NewHisTransReq(1, 20170414, "600000", 2000, 1)
	buf := new(bytes.Buffer)
	req.Write(buf)
	return buf, req
}

func _TestHisTransReq(t *testing.T) {
	fmt.Println("TestHisTransReq...")
	buf, req := BuildHisTransBuffer()

	start := time.Now().UnixNano()
	conn, err := net.Dial("tcp", HOST)
	chk(err)

	_, err = conn.Write(buf.Bytes())
	chk(err)

	err, buffer := entity.ReadResp(conn)
	chk(err)
	fmt.Println("time cost: ", time.Now().UnixNano() - start)

	parser := entity.NewHisTransParser(req, buffer)
	result := parser.Parse()

	fmt.Println("record count: ", len(result))
	for _, t := range result {
		fmt.Println(t)
	}
}

func BuildPeriodDataBuffer() (*bytes.Buffer, *entity.PeriodDataReq) {
	req := entity.NewPeriodDataReq(1, "600000", entity.PERIOD_DAY, 0, 0x118)
	buf := new(bytes.Buffer)
	req.Write(buf)
	return buf, req
}

func TestPeriodDataReq(t *testing.T) {
	fmt.Println("TestPeriodDataReq...")
	buf, req := BuildPeriodDataBuffer()

	start := time.Now().UnixNano()
	conn, err := net.Dial("tcp", HOST)
	chk(err)

	_, err = conn.Write(buf.Bytes())
	chk(err)

	err, buffer := entity.ReadResp(conn)
	chk(err)
	fmt.Println("time cost: ", time.Now().UnixNano() - start)

	parser := entity.NewPeriodDataParser(req, buffer)
	result := parser.Parse()
	fmt.Println(hex.EncodeToString(parser.Data))

	fmt.Println("record count: ", len(result))
	for _, t := range result {
		fmt.Println(t)
	}
}

func _TestReqData(t *testing.T) {
	reqHex := "0c57086401011c001c002d050000333030313732040001000100180100000000000000000000"
	reqData, _ := hex.DecodeString(reqHex)

	conn, err := net.Dial("tcp", HOST)
	chk(err)

	_, err = conn.Write(reqData)
	chk(err)

	err, buffer := entity.ReadResp(conn)
	chk(err)

	parser := entity.NewRespParser(buffer)
	parser.Parse()

	fmt.Println(hex.EncodeToString(parser.Data))
}

func Test(t *testing.T) {
}
