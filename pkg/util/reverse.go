package util

func ReverseMap[M ~map[K]V, K comparable, V comparable](m M) map[V]K {
	reversedMap := make(map[V]K)
	for key, value := range m {
		reversedMap[value] = key
	}
	return reversedMap
}
