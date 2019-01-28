package service

import (
	"encoding/json"
	"github.com/bytom/errors"

	"github.com/bufferserver/api/common"
	"github.com/bufferserver/util"
)

// Node can invoke the api which provide by the full node server
type Node struct {
	url string
}

// Node create a api client with target server
func NewNode(url string) *Node {
	return &Node{url: url}
}

type attachUtxo struct {
	Hash   string `json:"hash"`
	Asset  string `json:"asset"`
	Amount uint64 `json:"amount"`
}

func (n *Node) ListBlockCenterUTXOs(req *common.Display) ([]*attachUtxo, error) {
	url := "/api/v1/btm/q/list-utxos"
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}

	resp, err := n.request(url, payload)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var res []*attachUtxo
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

func (n *Node) request(path string, payload []byte) (interface{}, error) {
	resp := &Response{}
	if err := util.Post(n.url+path, payload, resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, errors.New(resp.Msg)
	}

	return resp.Result["data"], nil
}
