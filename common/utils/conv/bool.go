package conv

func B2I[T int | int8 | int16 | int32 | int64](b bool) T {
	if b {
		return T(1)
	} else {
		return T(0)
	}
}

func I2B[T int | int8 | int16 | int32 | int64](i T) bool {
	if i == 0 {
		return false
	} else {
		return true
	}
}
