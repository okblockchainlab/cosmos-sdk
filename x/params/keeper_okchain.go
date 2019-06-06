package params


func (k *Keeper) RegisterParamSet(paramSpace string, ps ...ParamSet) *Keeper {
	for _, ps := range ps {
		if ps != nil {
			// if _, ok := paramSets[ps.GetParamSpace()]; ok {
			// 	panic(fmt.Sprintf("<%s> already registered ", ps.GetParamSpace()))
			// }
			k.paramSets[paramSpace] = ps
		}
	}
	return k
}

// Get existing substore from keeper
func (k Keeper) GetParamSet(paramSpace string) (ParamSet, bool) {
	paramSet, ok := k.paramSets[paramSpace]
	if !ok {
		return nil, false
	}
	return paramSet, ok
}