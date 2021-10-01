package internal

// ValueOrDefaultString safe dereferencing of strings
func ValueOrDefaultString(s *string, defaultString ...string) string {
	def := ""
	if len(defaultString) > 0 {
		def = defaultString[0]
	}
	if s != nil {
		def = *s
	}
	return def
}

func Contains(list []string, value string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}
