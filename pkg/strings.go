package pkg

// EqualAny returns whether a string is equal to any of the given strings.
func EqualAny(a string, b ...string) bool {
	for _, s := range b {
		if a == s {
			return true
		}
	}
	return false
}
