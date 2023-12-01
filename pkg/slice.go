package pkg

func UnsetSlice[T any](s []T, i int) []T {
	if i > len(s)-1 {
		return s
	}
	if i < 0 {
		return s
	}
	return append(s[:i], s[i+1:]...)
}
