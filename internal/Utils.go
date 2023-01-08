package internal

func MapArray[K any, T any](data []K, mapper func(idx int, k *K) T) []T {
	res := make([]T, len(data))
	for idx, k := range data {
		res[idx] = mapper(idx, &k)
	}
	return res
}
