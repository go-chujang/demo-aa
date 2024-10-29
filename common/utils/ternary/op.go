package ternary

import "reflect"

func Cond[T any](cond bool, trueCase, falseCase T) T {
	if cond {
		return trueCase
	}
	return falseCase
}

func Fn[T any](cond bool, trueCase T, fn func() T) T {
	if cond {
		return trueCase
	}
	return fn()
}

func Err[T any](err error, nilCase, nonNilCase T) T {
	return Cond(err == nil, nilCase, nonNilCase)
}

func Default[T any](fn func(T) bool, value, defaultValue T) T {
	fn = Cond(fn != nil, fn, func(v T) bool {
		return !reflect.ValueOf(v).IsZero()
	})
	if fn(value) {
		return value
	}
	return defaultValue
}

func VArgs[T any](fn func(T) bool, defaultValue T, args ...T) T {
	fn = Cond(fn != nil, fn, func(v T) bool {
		return !reflect.ValueOf(v).IsZero()
	})
	for _, v := range args {
		if fn(v) {
			return v
		}
	}
	return defaultValue
}
