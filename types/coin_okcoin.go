package types

import "encoding/json"

// MarshalJSON marshals the coin
func (coin Coin) MarshalJSON() ([]byte, error) {
	type Alias Coin
	return json.Marshal(&struct {
		Denom  string `json:"denom"`
		Amount Dec    `json:"amount"`
	}{
		coin.Denom,
		NewDecFromIntWithPrec(coin.Amount, Precision),
	})
}

func (coin *Coin) UnmarshalJSON(data []byte) error {
	c := &struct {
		Denom  string `json:"denom"`
		Amount Dec    `json:"amount"`
	}{}
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}
	coin.Denom = c.Denom
	coin.Amount = NewIntFromBigInt(c.Amount.Int)
	return nil
}