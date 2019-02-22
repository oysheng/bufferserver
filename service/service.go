package service

import (
	"encoding/json"
	"github.com/bytom/errors"

	"github.com/bufferserver/api/common"
	"github.com/bufferserver/util"
)

// Service can invoke the api which provide by the server
type Service struct {
	url string
}

// NewService new a service with target server
func NewService(url string) *Service {
	return &Service{url: url}
}

type AttachUtxo struct {
	Hash   string `json:"hash"`
	Asset  string `json:"asset"`
	Amount uint64 `json:"amount"`
}

func (s *Service) ListBlockCenterUTXOs(req *common.Display) ([]*AttachUtxo, error) {
	url := "/api/v1/btm/q/list-utxos"
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}

	resp, err := s.request(url, payload)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var res []*AttachUtxo
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

type Response struct {
	Code   int                    `json:"code"`
	Msg    string                 `json:"msg"`
	Result map[string]interface{} `json:"result,omitempty"`
}

func (s *Service) request(path string, payload []byte) (interface{}, error) {
	resp := &Response{}
	if err := util.Post(s.url+path, payload, resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, errors.New(resp.Msg)
	}

	return resp.Result["data"], nil
}

type TransactionStatusResp struct {
	Height     int64 `json:"height"`
	StatusFail bool  `json:"status_fail"`
}

func (s *Service) GetTransactionStatus(TxID string) (*TransactionStatusResp, error) {
	url := s.url + "/transaction/" + TxID
	var resp TransactionStatusResp
	if err := util.Get(url, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
