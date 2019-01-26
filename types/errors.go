package types

import (
	"github.com/bytom/errors"
)

var (
	// general error
	ErrMissWallet       = errors.New("can't find the wallet")
	ErrMissXPub         = errors.New("can't find the xpub")
	ErrAssetID          = errors.New("asset id format error")
	ErrAddressFormat    = errors.New("address format error")
	ErrPlatformNotFound = errors.New("platform not found")
	ErrBannedIPOrWallet = errors.New("the ip or wallet is banned")

	// build payment error
	ErrBuildPayment           = errors.New("build payment fail")
	ErrEmptyBuildToAddr       = errors.New("empty build payment destination address")
	ErrEmptyBuildAmount       = errors.New("empty build amount")
	ErrInsufficientSpendUTXOs = errors.New("insufficient spend utxos")
	ErrInsufficientFeeUTXOs   = errors.New("insufficient fee utxos")
	ErrOutputAmoutTooLarge    = errors.New("output amount exceed maximum value 2^63")

	ErrBadBuildType = errors.New("bad build type")

	// submit tx error
	ErrFinalizeTx = errors.New("finalize tx fail")
	ErrSubmitTx   = errors.New("submit tx fail")
)
