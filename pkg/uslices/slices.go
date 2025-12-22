package uslices

// Unique 去重
func Unique[T comparable](s []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(s))

	for _, item := range s {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
