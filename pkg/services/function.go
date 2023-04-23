package services

func Filter[T any](slice []T, condition func(T) bool) ([]T, []T) {
	var filtered []T
	var differed []T
	for _, item := range slice {
		if condition(item) {
			filtered = append(filtered, item)
		} else {
			differed = append(differed, item)
		}
	}
	return filtered, differed
}

func ForEach[T any](slice []T, action func(T)) {
	for _, item := range slice {
		action(item)
	}
}

func Map[T1 any, T2 any](slice []T1, mapper func(T1) T2) []T2 {
	mapped := make([]T2, 0, len(slice))
	for _, item := range slice {
		mapped = append(mapped, mapper(item))
	}
	return mapped
}
