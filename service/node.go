package service

import (
	"encoding/json"
	"github.com/bytom/errors"

	"github.com/bufferserver/api"
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

func (n *Node) ListBlockCenterUTXOs(req *common.Display) ([]*api.ListUTXOsResp, error) {
	url := "/api/v1/btm/q/list-utxos"
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}

	res := []*api.ListUTXOsResp{}
	return res, n.request(url, payload, res)
}

type response struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	ErrDetail string          `json:"error_detail"`
}

func (n *Node) request(path string, payload []byte, respData interface{}) error {
	resp := &response{}
	if err := util.Post(n.url+path, payload, resp); err != nil {
		return err
	}

	if resp.Status != "success" {
		return errors.New(resp.ErrDetail)
	}

	return json.Unmarshal(resp.Data, respData)
}
