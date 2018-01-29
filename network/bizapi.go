package network

import (
	"fmt"
	"errors"
	"io/ioutil"
	"path/filepath"
	"os"
	"time"
	"github.com/stephenlyu/tds/datasource/tdx"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tds/date"
)

var blockExchangeMap = map[uint16]string{
	0: "SZ",
	1: "SH",
}

type BizApi struct {
	api *API

	workDir string
}

func CreateBizApi(host string) (error, *BizApi) {
	result := &BizApi{workDir: "temp"}
	err, api := CreateAPI(fmt.Sprintf("%s:7709", host))
	if err != nil {
		return err, nil
	}

	result.api = api

	return nil, result
}

func (this *BizApi) Cleanup() {
	if this.api != nil {
		this.api.Cleanup()
		this.api = nil
	}
}

func (this *BizApi) SetTimeOut(timeout int) {
	this.api.SetTimeOut(timeout)
}

func (this *BizApi) SetWorkDir(dir string) {
	this.workDir = dir
}

func (this *BizApi) getStockCodesByBlock(block uint16) (error, []string) {
	exchange, ok := blockExchangeMap[block]
	if !ok {
		return nil, nil
	}

	outputDir := filepath.Join(this.workDir, "T0002/hq_cache")
	zhbFile := "zhb.zip"

	zhbFilePath := filepath.Join(outputDir, zhbFile)
	stats, err := os.Stat(zhbFilePath)

	today := date.GetTodayString()

	if os.IsNotExist(err) || date.ToDayString(stats.ModTime()) < today {
		err := this.DownloadFile(zhbFile, outputDir)
		if err != nil {
			return err, nil
		}
		err = util.UnzipFile(zhbFilePath, outputDir)
		if err != nil {
			return err, nil
		}
	}

	ds := tdxdatasource.NewDataSource(this.workDir, true)
	return nil, ds.GetStockCodes(exchange)
}

func (this *BizApi) GetSZStockCodes() (error, []string) {
	return this.getStockCodesByBlock(BLOCK_SZ_A)
}

func (this *BizApi) GetSHStockCodes() (error, []string) {
	return this.getStockCodesByBlock(BLOCK_SH_A)
}

func (this *BizApi) GetAStockCodes() (error, []string) {
	result := []string{}

	err, codes := this.GetSZStockCodes()
	if err != nil {
		return err, nil
	}

	result = append(result, codes...)

	err, codes = this.GetSHStockCodes()
	if err != nil {
		return err, nil
	}

	result = append(result, codes...)
	return nil, result
}

func (this *BizApi) GetInfoEx(codes []string) (error, map[string][]*InfoExItem) {
	result := map[string][]*InfoExItem{}

	n := 20
	for i := 0; i < len(codes); i += n {
		end := i + n
		if end > len(codes) {
			end = len(codes)
		}
		subCodes := codes[i:end]
		err, infoEx := this.api.GetInfoEx(subCodes)
		if err != nil {
			return err, nil
		}

		for k, v := range infoEx {
			result[k] = v
		}
	}

	return nil, result
}

func (this *BizApi) GetFinance(codes []string) (error, map[string]*Finance) {
	result := map[string]*Finance{}

	n := 100
	for i := 0; i < len(codes); i += n {
		end := i + n
		if end > len(codes) {
			end = len(codes)
		}
		subCodes := codes[i:end]
		err, finances := this.api.GetFinance(subCodes)
		if err != nil {
			return err, nil
		}

		for k, v := range finances {
			result[k] = v
		}
	}

	return nil, result
}

func (this *BizApi) GetLatestMinuteData(code string, offset int, count int) (error, []*Record) {
	result := []*Record{}

	n := 0

	for n < count {
		c := 280
		if c > count - n {
			c = count - n
		}

		err, data := this.api.GetMinuteData(code, uint16(offset + n), uint16(c))
		if err != nil {
			return err, nil
		}

		if len(data) == 0 {
			break
		}

		result = append(data, result...)
		n += len(data)
	}

	return nil, result
}

func (this *BizApi) GetLatestDayData(code string, count int) (error, []*Record) {
	result := []*Record{}

	n := 0

	for n < count {
		c := 280
		if c > count - n {
			c = count - n
		}

		err, data := this.api.GetDayData(code, uint16(n), uint16(c))
		if err != nil {
			return err, nil
		}

		if len(data) == 0 {
			break
		}

		result = append(data, result...)
		n += len(data)
	}

	return nil, result
}

func (this *BizApi) DownloadFile(fileName string, outputDir string) error {
	err, length := this.api.GetFileLength(fileName)
	if err != nil {
		return err
	}

	fileData := make([]byte, length)

	var offset uint32 = 0
	var count uint32 = 30000

	var getPacket = func() (error error, packetLength uint32, data []byte) {
		retryTimes := 0
		for retryTimes < 3 {
			err, packetLength, data = this.api.GetFileData(fileName, offset, count)
			if err == nil {
				return
			}
			time.Sleep(time.Millisecond * 500)
			retryTimes++
		}
		return
	}

	for offset < length {
		err, packetLength, data := getPacket()
		if err != nil {
			return err
		}
		if packetLength != uint32(len(data)) {
			return errors.New("bad data")
		}

		copy(fileData[offset:offset + packetLength], data[:])

		offset += count
	}

	filePath := filepath.Join(outputDir, fileName)
	os.MkdirAll(filepath.Dir(filePath), 0777)
	return ioutil.WriteFile(filePath, fileData, 0666)
}
