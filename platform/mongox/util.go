package mongox

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func StructOrMap2BasicUpdateSet(out *bson.M, in interface{}, prefix string, adds ...bson.M) bool {
	if out == nil {
		return false
	}

	rv := reflect.ValueOf(in)
	if rv.IsZero() && adds == nil {
		return false
	}
	if prefix != "" && !strings.HasSuffix(prefix, ".") {
		prefix += "."
	}
	if adds != nil {
		for key, value := range adds[0] {
			(*out)[prefix+key] = value
		}
	}

	switch rv.Type().Kind() {
	case reflect.Struct:
		for i := 0; i < rv.Type().NumField(); i++ {
			value := rv.Field(i)
			if value.IsZero() {
				continue
			}

			field := rv.Type().Field(i)
			tag := field.Tag.Get("json")
			if tag == "-" {
				continue
			}
			if tag == "" {
				tag = field.Name
			}

			tag = strings.Split(tag, ",")[0] // trimSuffix 'omitempty' or etc
			tag = prefix + tag

			switch value.Kind() {
			case reflect.Array,
				reflect.Slice,
				reflect.Chan,
				reflect.Func:
				// unsupported
				continue
			case reflect.Struct, reflect.Map:
				StructOrMap2BasicUpdateSet(out, value.Interface(), tag)
				continue
			}
			(*out)[tag] = value.Interface()
		}
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			if key.Kind() != reflect.String {
				continue
			}

			tag := prefix + key.String()
			value := rv.MapIndex(key)
			if value.IsZero() {
				continue
			}

			switch value.Kind() {
			case reflect.Array,
				reflect.Slice,
				reflect.Chan,
				reflect.Func:
				// unsupported
				continue
			case reflect.Struct, reflect.Map:
				StructOrMap2BasicUpdateSet(out, value.Interface(), tag)
				continue
			}
			(*out)[tag] = value.Interface()
		}
	default:
		return false
	}
	return len((*out)) != 0
}
