package helper

func AssignStringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}
