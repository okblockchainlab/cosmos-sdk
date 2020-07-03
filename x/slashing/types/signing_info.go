package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValStatus int32

const (
	Created             ValStatus = 0x00
	Destroying          ValStatus = 0x01
	Destroyed           ValStatus = 0x02
	ValStatusCreated              = "Created"
	ValStatusDestroying           = "Destroying"
	ValStatusDestroyed            = "Destroyed"
)

// NewValidatorSigningInfo creates a new ValidatorSigningInfo instance
func NewValidatorSigningInfo(
	condAddr sdk.ConsAddress, startHeight, indexOffset int64,
	jailedUntil time.Time, tombstoned bool, missedBlocksCounter int64, validatorStatus ValStatus,
) ValidatorSigningInfo {

	return ValidatorSigningInfo{
		Address:             condAddr,
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
		ValidatorStatus:     validatorStatus,
	}
}

// String implements the stringer interface for ValidatorSigningInfo
func (i ValidatorSigningInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
  Address:               %s
  Start Height:          %d
  Index Offset:          %d
  Jailed Until:          %v
  Tombstoned:            %t
  Missed Blocks Counter: %d`,
		i.Address, i.StartHeight, i.IndexOffset, i.JailedUntil,
		i.Tombstoned, i.MissedBlocksCounter)
}

// unmarshal a validator signing info from a store value
func UnmarshalValSigningInfo(cdc codec.Marshaler, value []byte) (signingInfo ValidatorSigningInfo, err error) {
	err = cdc.UnmarshalBinaryBare(value, &signingInfo)
	return signingInfo, err
}
