package configor

func defaultValue(value, defaultVal string) string {
	if value == "" {
		return defaultVal
	}
	return value
}

// defaultValuePos returns value at position
func defaultValuePos(values []string, idx int, defaultVal string) string {
	if values == nil || idx >= len(values) {
		return defaultVal
	}

	return values[idx]
}
