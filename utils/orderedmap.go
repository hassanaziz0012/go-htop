package utils

type OrderedMap[T any] struct {
	Keys []string
	Data map[string]T
}

func NewOrderedMap[T any](keys []string, data map[string]T) *OrderedMap[T] {
	return &OrderedMap[T]{
		keys,
		data,
	}
}
