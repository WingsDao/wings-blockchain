package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
	"wings-blockchain/x/poa/types"
	"wings-blockchain/helpers"
)

// Type for codec
const (
	MsgAddValidatorType = types.ModuleName + "/add-validator"
)

// Message for adding validator
type MsgAddValidator struct {
	Address		sdk.AccAddress
	EthAddress 	string
	Sender		sdk.AccAddress
}

// Create new 'add validator' message
func NewMsgAddValidator(address sdk.AccAddress, ethAddress string, sender sdk.AccAddress) MsgAddValidator {
	return MsgAddValidator{
		Address: 	address,
		EthAddress: ethAddress,
		Sender:		sender,
	}
}

// Message route
func (msg MsgAddValidator) Route() string {
	return types.DefaultRoute
}

// Message type
func (msg MsgAddValidator) Type() string {
	return "add_validator"
}

// Validate basic for add validator msg
func (msg MsgAddValidator) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	if len(msg.EthAddress) == 0 {
		return sdk.ErrUnknownRequest("Wrong Ethereum address for validator")
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	if !helpers.IsEthereumAddress(msg.EthAddress) {
		return types.ErrWrongEthereumAddress(msg.EthAddress)
	}

	return nil
}

// Get signature bytes
func (msg MsgAddValidator) GetSignBytes() []byte {
	b, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// Get signers addresses
func (msg MsgAddValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}