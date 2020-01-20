package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeEmptyContractCode = 101

	CodeCantConnectVM = 201
	CodeErrDuringExec = 202

	CodeErrWrongModuleAddress = 301
	CodeErrModuleExists       = 302
)

func ErrEmptyContract() sdk.Error {
	return sdk.NewError(Codespace, CodeEmptyContractCode, "contract code is empty, please fill field with compiled contract bytes")
}

func ErrCantConnectVM(msg string) sdk.Error {
	return sdk.NewError(Codespace, CodeCantConnectVM, "cant connect to vm instance: %s", msg)
}

func ErrDuringVMExec(msg string) sdk.Error {
	return sdk.NewError(Codespace, CodeErrDuringExec, "cant execute contract: %s", msg)
}

func ErrWrongModuleAddress(expected, real sdk.AccAddress) sdk.Error {
	return sdk.NewError(Codespace, CodeErrWrongModuleAddress, "wrong module owner %s address, expected %s", expected, real)
}

func ErrModuleExists(address sdk.AccAddress, path []byte) sdk.Error {
	return sdk.NewError(Codespace, CodeErrModuleExists, "module %s already exists for account %s", hex.EncodeToString(path), address)
}
