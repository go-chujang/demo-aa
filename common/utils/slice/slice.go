package slice

func WithDefault[T any](defaultValue T, size int) []T {
	slice := make([]T, size)
	for i := 0; i < size; i++ {
		slice[i] = defaultValue
	}
	return slice
}
