package sliceutils

// ConvertFuncWithSkip ...
func ConvertFuncWithSkip[From, To any](slice []From, convertFunc func(elem From) (To, bool)) []To {
	if len(slice) == 0 {
		return nil
	}

	result := make([]To, 0, len(slice))

	for _, elem := range slice {
		resElem, skip := convertFunc(elem)
		if skip {
			continue
		}

		result = append(result, resElem)
	}

	return result
}

// ConvertFunc ...
func ConvertFunc[From, To any](slice []From, convertFunc func(elem From) To) []To {
	return ConvertFuncWithSkip(
		slice,
		func(elem From) (To, bool) {
			return convertFunc(elem), false
		},
	)
}
