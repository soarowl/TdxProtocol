package network

import (
	"sync"
	"gopkg.in/fatih/pool.v2"
	"net"
	"bytes"
)

type API struct {
	seqId			uint32
	lock    		sync.Mutex

	pool 			pool.Pool
}

func CreateAPI(host string) (error, *API) {
	api := &API {}
	err := api.Initialize(host)
	if err != nil {
		return err, nil
	}

	return nil, api
}

func (this *API) Initialize(host string) error {
	factory := func() (net.Conn, error) {
		return net.Dial("tcp", host)
	}

	p, err := pool.NewChannelPool(1, 1, factory)
	if err != nil {
		return err
	}

	this.pool = p

	return nil
}

func (this *API) Cleanup() error {
	if this.pool != nil {
		this.pool.Close()
		this.pool = nil
	}
	return nil
}

func (this *API) nextSeqId() uint32 {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.seqId++
	return this.seqId
}

func (this *API) markConnUnusable(conn interface{}) {
	if poolConn, ok := conn.(pool.PoolConn); ok {
		poolConn.MarkUnusable()
	}
}

func (this *API) sendReq(data []byte) (error, []byte) {
	conn, err := this.pool.Get()
	if err != nil {
		return err, nil
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		this.markConnUnusable(conn)
		return nil, nil
	}

	err, respData := ReadResp(conn)
	if err != nil {
		this.markConnUnusable(conn)
		return err, nil
	}

	return err, respData
}

func (this *API) GetStockList(block uint16, offset uint16, count uint16) (error, int, map[string]*Bid) {
	req := NewStockListReq(this.nextSeqId(), block, offset, count)
	buf := new(bytes.Buffer)
	req.Write(buf)

	err, respData := this.sendReq(buf.Bytes())
	if err != nil {
		return err, 0, nil
	}

	parser := NewStockListParser(req, respData)
	result := parser.Parse()

	return nil, int(parser.Total), result
}

func (this *API) GetInfoEx(codes []string) (error, map[string][]*InfoExItem) {
	req := NewInfoExReq(this.nextSeqId())
	for _, code := range codes {
		req.AddCode(code)
	}
	buf := new(bytes.Buffer)
	req.Write(buf)

	err, respData := this.sendReq(buf.Bytes())
	if err != nil {
		return err, nil
	}

	parser := NewInfoExParser(req, respData)
	result := parser.Parse()

	return nil, result
}

func (this *API) GetPeriodData(code string, period, offset, count uint16) (error, []*Record) {
	req := NewPeriodDataReq(this.nextSeqId(), code, period, offset, count)
	buf := new(bytes.Buffer)
	req.Write(buf)

	err, respData := this.sendReq(buf.Bytes())
	if err != nil {
		return err, nil
	}

	parser := NewPeriodDataParser(req, respData)
	result := parser.Parse()

	return nil, result
}

func (this *API) GetMinuteData(code string, offset, count uint16) (error, []*Record) {
	return this.GetPeriodData(code, PERIOD_MINUTE, offset, count)
}

func (this *API) GetDayData(code string, offset, count uint16) (error, []*Record) {
	return this.GetPeriodData(code, PERIOD_DAY, offset, count)
}
