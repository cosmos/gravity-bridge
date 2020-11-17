package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

func (s GenesisState) ValidateBasic() error {
	if err := s.Params.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}
	return nil
}
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: &Params{},
	}
}
