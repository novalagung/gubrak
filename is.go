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

// IsArray will return true when type of the data is array/slice
func IsArray(data interface{}) bool {
	return _typeIs(data,
		reflect.Array,
		reflect.Slice,
	)
}

// IsBool will return true when type of the data is boolean
func IsBool(data interface{}) bool {
	return _typeIs(data,
		reflect.Bool,
	)
}

// IsChannel will return true when type of the data is channel
func IsChannel(data interface{}) bool {
	return _typeIs(data,
		reflect.Chan,
	)
}

// IsDate will return true when type of the data is time.Time
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
	} else if value, ok := data.(string); ok {
		return value == ""
	} else if value, ok := data.(bool); ok {
		return value == true
	} else if value, ok := data.(float32); ok {
		return value == 0
	} else if value, ok := data.(float64); ok {
		return value == 0
	} else if value, ok := data.(int); ok {
		return value == 0
	} else if value, ok := data.(int8); ok {
		return value == 0
	} else if value, ok := data.(int16); ok {
		return value == 0
	} else if value, ok := data.(int32); ok {
		return value == 0
	} else if value, ok := data.(int64); ok {
		return value == 0
	} else if value, ok := data.(uint); ok {
		return value == 0
	} else if value, ok := data.(uint8); ok {
		return value == 0
	} else if value, ok := data.(uint16); ok {
		return value == 0
	} else if value, ok := data.(uint32); ok {
		return value == 0
	} else if value, ok := data.(uint64); ok {
		return value == 0
	} else if IsNil(data) {
		return true
	} else if IsMap(data) || IsArray(data) {
		return valueOf(data).Len() == 0
	}

	return false
}

// IsEmptyString will return true when type of the data is string and it's empty
func IsEmptyString(data interface{}) bool {
	if data == nil {
		return true
	}

	if value, ok := data.(string); ok {
		return value == ""
	}

	return false
}

// IsFloat will return true when type of the data is floating number
func IsFloat(data interface{}) bool {
	return _typeIs(data,
		reflect.Float32,
		reflect.Float64,
	)
}

// IsFunction will return true when type of the data is closure/function
func IsFunction(data interface{}) bool {
	return _typeIs(data,
		reflect.Func,
	)
}

// IsInt will return true when type of the data is numeric integer
func IsInt(data interface{}) bool {
	return _typeIs(data,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
	)
}

// IsMap will return true when type of the data is hash map
func IsMap(data interface{}) bool {
	return _typeIs(data,
		reflect.Map,
	)
}

// IsNil will return true when type of the data is nil
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

// IsNumeric will return true when type of the data is numeric (float, uint, int)
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

// IsPointer will return true when type of the data is pointer
func IsPointer(data interface{}) bool {
	return _typeIs(data,
		reflect.Ptr,
	)
}

// IsSlice is alias of IsArray()
func IsSlice(data interface{}) bool {
	return IsArray(data)
}

// IsStructObject will return true when type of the data is object from struct
func IsStructObject(data interface{}) bool {
	return _typeIs(data,
		reflect.Struct,
	)
}

// IsString will return true when type of the data is string
func IsString(data interface{}) bool {
	return _typeIs(data,
		reflect.String,
	)
}

// IsTrue will return true when type of the data is bool, and the value is true
func IsTrue(data interface{}) bool {
	if data == nil {
		return true
	}

	if value, ok := data.(bool); ok {
		return value == true
	}

	return false
}

// IsUint will return true when type of the data is uint
func IsUint(data interface{}) bool {
	return _typeIs(data,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
	)
}

// IsZeroNumber will return true when type of the data is numeric and it's has 0 value
func IsZeroNumber(data interface{}) bool {
	if data == nil {
		return true
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

	return false
}
