package network

import "fmt"

type BizApi struct {
	api *API
}

func CreateBizApi(host string) (error, *BizApi) {
	result := &BizApi{}
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

func (this *BizApi) getStockCodesByBlock(block uint16) (error, []string) {
	err, total, _ := this.api.GetStockList(block, 0, 1)
	if err != nil {
		return err, nil
	}

	result := make([]string, total)
	count := 0

	for count < total {
		err, _, bids := this.api.GetStockList(block, uint16(count), 80)
		if err != nil {
			return err, nil
		}

		for k, _ := range bids {
			result[count] = k
			count++
		}
	}

	return nil, result
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

func (this *BizApi) GetLatestMinuteData(code string, count int) (error, []*Record) {
	result := []*Record{}

	n := 0

	for n < count {
		c := 280
		if c > count - n {
			c = count - n
		}

		err, data := this.api.GetMinuteData(code, uint16(n), uint16(c))
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
