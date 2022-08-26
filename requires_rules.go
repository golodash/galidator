package galidator

// Returns false if one of the fields values is not empty, nil or zero and input is empty, nil or zero
func whenExistOneRequire(fields ...string) func(interface{}) func(interface{}) bool {
	return func(all interface{}) func(interface{}) bool {
		fieldsValues := getValues(all, fields...)
		return func(input interface{}) bool {
			for _, element := range fieldsValues {
				if !isEmptyNilZero(element) {
					if isEmptyNilZero(input) {
						return false
					}
				}
			}
			return true
		}
	}
}
