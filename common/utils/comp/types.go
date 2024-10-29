package comp

import "reflect"

func TypesEqual(a, b any) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}
