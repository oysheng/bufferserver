package common

import (
	"github.com/bytom/errors"

	"github.com/blockcenter/types"
)

//FormatErrResp format error response
func formatErrResp(err error) Response {
	// default error response
	response := Response{
		Code: 300,
		Msg:  "request error",
	}

	root := errors.Root(err)
	if errCode, ok := respErrFormatter[root]; ok {
		response.Code = errCode
		response.Msg = root.Error()
	}
	return response
}

var respErrFormatter = map[error]int{
	// general error
	types.ErrAssetID:          400,
	types.ErrAddressFormat:    401,
	types.ErrPlatformNotFound: 402,
	types.ErrBannedIPOrWallet: 403,
	types.ErrMissXPub:         404,

	// build payment error
	types.ErrBuildPayment:           500,
	types.ErrEmptyBuildAmount:       501,
	types.ErrEmptyBuildToAddr:       502,
	types.ErrInsufficientSpendUTXOs: 503,
	types.ErrInsufficientFeeUTXOs:   504,
	types.ErrOutputAmoutTooLarge:    505,

	// submit tx error
	types.ErrFinalizeTx: 600,
	types.ErrSubmitTx:   601,
}
