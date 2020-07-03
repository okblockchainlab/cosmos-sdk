package types

// GenesisState - minter state
type GenesisState struct {
	Minter MinterCustom `json:"minter_custom" yaml:"minter_custom"` // minter object
	Params Params       `json:"params" yaml:"params"`               // inflation params
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter MinterCustom, params Params) GenesisState {
	return GenesisState{
		Minter: minter,
		Params: params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Minter: DefaultInitialMinterCustom(),
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return ValidateMinterCustom(data.Minter)
}
