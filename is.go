package gubrak

import (
	"reflect"
	"time"
)

func _typeIs(data interface{}, types ...reflect.Kind) bool {

	if dataKind, ok := data.(reflect.Kind); ok {
		for _, each := range types {
			if dataKind == each {
				return true
			}
		}
	}

	var err error

	result := func(err *error) bool {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return false
		}

		dataKind := typeOf(data).Kind()

		for _, each := range types {
			if dataKind == each {
				return true
			}
		}

		return false
	}(&err)

	if err != nil {
		return false
	}

	return result
}

func IsArray(data interface{}) bool {
	return _typeIs(data,
		reflect.Array,
		reflect.Slice,
	)
}

func IsBool(data interface{}) bool {
	return _typeIs(data,
		reflect.Bool,
	)
}

func IsChannel(data interface{}) bool {
	return _typeIs(data,
		reflect.Chan,
	)
}

func IsDate(data interface{}) bool {
	if _, ok := data.(time.Time); ok {
		return true
	}

	return false
}

// IsEmpty will return false to any null-able data which value is nil (chan, func, interface, map, pointer, slice), will also return false when the value is default value of it's data type (false for bool, "" for string, 0 for numeric value), and will return false if the value is slice or map and the length is 0
func IsEmpty(data interface{}) bool {
	if data == nil {
		return true
	}

	if value, ok := data.(string); ok {
		return value == ""
	}

	if value, ok := data.(bool); ok {
		return value == true
	}

	if value, ok := data.(float32); ok {
		return value == 0
	}
	if value, ok := data.(float64); ok {
		return value == 0
	}

	if value, ok := data.(int); ok {
		return value == 0
	}
	if value, ok := data.(int8); ok {
		return value == 0
	}
	if value, ok := data.(int16); ok {
		return value == 0
	}
	if value, ok := data.(int32); ok {
		return value == 0
	}
	if value, ok := data.(int64); ok {
		return value == 0
	}

	if value, ok := data.(uint); ok {
		return value == 0
	}
	if value, ok := data.(uint8); ok {
		return value == 0
	}
	if value, ok := data.(uint16); ok {
		return value == 0
	}
	if value, ok := data.(uint32); ok {
		return value == 0
	}
	if value, ok := data.(uint64); ok {
		return value == 0
	}

	if IsNil(data) {
		return true
	}

	if size, _ := Size(data); size == 0 {
		return true
	}

	return false
}

func IsFloat(data interface{}) bool {
	return _typeIs(data,
		reflect.Float32,
		reflect.Float64,
	)
}

func IsFunction(data interface{}) bool {
	return _typeIs(data,
		reflect.Func,
	)
}

func IsInt(data interface{}) bool {
	return _typeIs(data,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
	)
}

func IsMap(data interface{}) bool {
	return _typeIs(data,
		reflect.Map,
	)
}

func IsNil(data interface{}) bool {
	if data == nil {
		return true
	}

	isNillable := _typeIs(data,
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice,
		reflect.Array,
	)
	if isNillable {
		return valueOf(data).IsNil()
	}

	return false
}

func IsNumeric(data interface{}) bool {
	return _typeIs(data,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Float32,
		reflect.Float64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
	)
}

func IsPointer(data interface{}) bool {
	return _typeIs(data,
		reflect.Ptr,
	)
}

func IsString(data interface{}) bool {
	return _typeIs(data,
		reflect.String,
	)
}

func IsObject(data interface{}) bool {
	return _typeIs(data,
		reflect.Struct,
	)
}

func IsUint(data interface{}) bool {
	return _typeIs(data,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
	)
}

// isArguments
// isArrayBuffer
// isArrayLike
// isArrayLikeObject
// isBuffer
// isDate
// isElement
// isEqual
// isEqualWith
// isError
// isFinite
// isFunction
// isLength
// isMatch
// isMatchWith
// isNaN
// isNative
// isObjectLike
// isPlainObject
// isRegExp
// isSafeInteger
// isSet
// isSymbol
// isTypedArray
// isUndefined
// isWeakMap
// isWeakSet
