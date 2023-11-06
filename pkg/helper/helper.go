package helper

func AssignValue(value interface{}) interface{} {
	if value == nil {
		return ""
	}
	return value
}
