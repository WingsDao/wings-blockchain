package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"wings-blockchain/x/vm/internal/types/vm_grpc"
)

const (
	EventKeyDiscard = "discard"
	EventKeyKeep    = "keep"
)

// New event with keep status.
func NewEventKeep() sdk.Event {
	return sdk.NewEvent(
		EventKeyKeep,
	)
}

// New event with discard status.
func NewEventDiscard(errorStatus *vm_grpc.VMErrorStatus) sdk.Event {
	attributes := make([]sdk.Attribute, 0)

	if errorStatus != nil {
		attributes = append(attributes, sdk.NewAttribute("major_status", string(errorStatus.MajorStatus)))
		attributes = append(attributes, sdk.NewAttribute("sub_status", string(errorStatus.SubStatus)))
		attributes = append(attributes, sdk.NewAttribute("message", errorStatus.Message))
	}

	return sdk.NewEvent(
		EventKeyDiscard,
		attributes...,
	)
}

// Parse VM event to standard SDK event.
// In case of event data equal "struct" we don't process struct, and just keep bytes, as for any other type.
func NewEventFromVM(event *vm_grpc.VMEvent) sdk.Event {
	return sdk.NewEvent(
		string(event.Key),
		sdk.NewAttribute("sequence_number", string(event.SequenceNumber)),
		sdk.NewAttribute("type", VMTypeToStringPanic(event.Type.Tag)),
		sdk.NewAttribute("event_data", hex.EncodeToString(event.EventData)),
	)
}
