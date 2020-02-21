// Basic constants and function to work with types.
package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/WingsDao/wings-blockchain/x/vm/internal/types/vm_grpc"
	"github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "vm"

	StoreKey  = ModuleName
	RouterKey = ModuleName

	Codespace         types.CodespaceType = ModuleName
	DefaultParamspace                     = ModuleName

	VmAddressLength = 32
	VmGasPrice      = 1
	VmUnknowTagType = -1
	zeroBytes = 12
)

// VM related variables.
var (
	KeyGenesisInitialized = []byte("gen_init")
	KeyDelimiter          = []byte(":")
	VMKey                 = []byte("vm")
)

// Type of Move contract (bytes).
type Contract []byte

// Convert bech32 to libra hex.
func Bech32ToLibra(acc types.AccAddress) string {
	prefix := types.GetConfig().GetBech32AccountAddrPrefix()
	zeros := make([]byte, zeroBytes-len(prefix))

	bytes := make([]byte, 0)
	bytes = append(bytes, []byte(prefix)...)
	bytes = append(bytes, zeros...)
	bytes = append(bytes, acc...)

	return hex.EncodeToString(bytes)
}

// Make path for storage from VMAccessPath.
func MakePathKey(path vm_grpc.VMAccessPath) []byte {
	return bytes.Join(
		[][]byte{
			VMKey,
			path.Address,
			path.Path,
		},
		KeyDelimiter,
	)
}

// Convert VMAccessPath to hex string
func PathToHex(path vm_grpc.VMAccessPath) string {
	return fmt.Sprintf("Access path: \n"+
		"\tAddress: %s\n"+
		"\tPath:    %s\n"+
		"\tKey:     %s\n", hex.EncodeToString(path.Address), hex.EncodeToString(path.Path), hex.EncodeToString(MakePathKey(path)))
}

// Get TypeTag by string TypeTag representation.
func GetVMTypeByString(typeTag string) (vm_grpc.VMTypeTag, error) {
	if val, ok := vm_grpc.VMTypeTag_value[typeTag]; !ok {
		return VmUnknowTagType, fmt.Errorf("can't find tag type %s, check correctness of type value", typeTag)
	} else {
		return vm_grpc.VMTypeTag(val), nil
	}
}

// Convert TypeTag to string representation.
func VMTypeToString(tag vm_grpc.VMTypeTag) (string, error) {
	if val, ok := vm_grpc.VMTypeTag_name[int32(tag)]; !ok {
		return "", fmt.Errorf("can't find string representation of type %d, check correctness of type value", tag)
	} else {
		return val, nil
	}
}

// Convert TypeTag to string representation with panic.
func VMTypeToStringPanic(tag vm_grpc.VMTypeTag) string {
	if val, ok := vm_grpc.VMTypeTag_name[int32(tag)]; !ok {
		panic(fmt.Errorf("can't find string representation of type %d, check correctness of type value", tag))
	} else {
		return val
	}
}
