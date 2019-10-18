package service

import (
	"encoding/json"

	"github.com/bufferserver/types"
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
	page := "?limit=1000&start=0"
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}

	resp, err := s.request(url+page, payload)
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

type GetTransactionReq struct {
	TxID string `json:"tx_id"`
}

func (s *Service) GetTransaction(req *GetTransactionReq) (*types.Tx, error) {
	urlPath := "/api/v1/btm/merchant/get-transaction"
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}

	resp, err := s.request(urlPath, payload)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var res *types.Tx
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}
