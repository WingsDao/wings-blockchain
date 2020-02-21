// Message to remove validator described.
package msgs

import (
	"encoding/json"

	"github.com/WingsDao/wings-blockchain/x/poa/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Type for codec
const (
	MsgRemoveValidatorType = types.ModuleName + "/remove-validator"
)

// Message for remove validator
type MsgRemoveValidator struct {
	Address sdk.AccAddress `json:"address"`
	Sender  sdk.AccAddress `json:"sender"`
}

// Creating 'remove validator' message
func NewMsgRemoveValidator(address sdk.AccAddress, sender sdk.AccAddress) MsgRemoveValidator {
	return MsgRemoveValidator{
		Address: address,
		Sender:  sender,
	}
}

// Message route
func (msg MsgRemoveValidator) Route() string {
	return types.RouterKey
}

// Message type
func (msg MsgRemoveValidator) Type() string {
	return "remove_validator"
}

// Validate basic for remove validator message
func (msg MsgRemoveValidator) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	return nil
}

// Get bytes to sign from message data
func (msg MsgRemoveValidator) GetSignBytes() []byte {
	b, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// Get signers addresses
func (msg MsgRemoveValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
