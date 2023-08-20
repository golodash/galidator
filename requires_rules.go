package galidator

// Returns false if one of the fields values is not empty, nil or zero and input is empty, nil or zero
//
// False means it is required and all validators have to check
func whenExistOneRequireRule(fields ...string) func(interface{}) func(interface{}) bool {
	return func(all interface{}) func(interface{}) bool {
		fieldsValues := getValues(all, fields...)
		return func(input interface{}) bool {
			inputIsNil := isEmptyNilZero(input)
			for _, element := range fieldsValues {
				if !isEmptyNilZero(element) && inputIsNil {
					return false
				}
			}
			return true
		}
	}
}

// Returns false if all of the fields values are not empty, nil or zero and input is empty, nil or zero
//
// False means it is required and all validators have to check
func whenExistAllRequireRule(fields ...string) func(interface{}) func(interface{}) bool {
	return func(all interface{}) func(interface{}) bool {
		fieldsValues := getValues(all, fields...)
		return func(input interface{}) bool {
			for _, element := range fieldsValues {
				if isEmptyNilZero(element) {
					return true
				}
			}
			return !isEmptyNilZero(input)
		}
	}
}

// Returns false if one of the fields values is empty, nil or zero and input is empty, nil or zero
//
// False means it is required and all validators have to check
func whenNotExistOneRequireRule(fields ...string) func(interface{}) func(interface{}) bool {
	return func(all interface{}) func(interface{}) bool {
		fieldsValues := getValues(all, fields...)
		return func(input interface{}) bool {
			inputIsNil := isEmptyNilZero(input)
			for _, element := range fieldsValues {
				if isEmptyNilZero(element) && inputIsNil {
					return false
				}
			}
			return true
		}
	}
}

// Returns false if all of the fields values are empty, nil or zero and input is empty, nil or zero
//
// False means it is required and all validators have to check
func whenNotExistAllRequireRule(fields ...string) func(interface{}) func(interface{}) bool {
	return func(all interface{}) func(interface{}) bool {
		fieldsValues := getValues(all, fields...)
		return func(input interface{}) bool {
			for _, element := range fieldsValues {
				if !isEmptyNilZero(element) {
					return true
				}
			}
			return !isEmptyNilZero(input)
		}
	}
}
