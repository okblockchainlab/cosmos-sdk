package types

// IsAllGTE returns true iff for every denom in coins, the denom is present at
// an equal or greater amount in coinsB.
// TODO: Remove once unsigned integers are used.
func (coins DecCoins) IsAllGTE(coinsB DecCoins) bool {
	diff, _ := coins.SafeSub(coinsB)
	if len(diff) == 0 {
		return true
	}

	return !diff.IsAnyNegative()
}

