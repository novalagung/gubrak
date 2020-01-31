package gubrak

import (
	"math"
	"reflect"
	"time"
)

func typeIs(data interface{}, types ...reflect.Kind) bool {
	valueOfData := reflect.ValueOf(data)
	for _, tipe := range types {
		if tipe == valueOfData.Kind() {
			return true
		}
	}

	return false
}

// IsSlice is alias of IsSlice()
func IsSlice(data interface{}) bool {
	return typeIs(data, reflect.Slice)
}

// IsArray is alias of IsArray()
func IsArray(data interface{}) bool {
	return typeIs(data, reflect.Array)
}

// IsSliceOrArray will return true when type of the data is array/slice
func IsSliceOrArray(data interface{}) bool {
	return IsSlice(data) || IsArray(data)
}

var IsArrayOrSlice = IsSliceOrArray

// IsBool will return true when type of the data is boolean
func IsBool(data interface{}) bool {
	return typeIs(data,
		reflect.Bool,
	)
}

// IsChannel will return true when type of the data is channel
func IsChannel(data interface{}) bool {
	return typeIs(data, reflect.Chan)
}

// IsDate will return true when type of the data is time.Time
func IsDate(data interface{}) bool {
	if _, ok := data.(time.Time); ok {
		return true
	}

	return false
}

// IsString will return true when type of the data is string
func IsString(data interface{}) bool {
	return typeIs(data, reflect.String)
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
	return typeIs(data,
		reflect.Float32,
		reflect.Float64,
	)
}

// IsFunction will return true when type of the data is closure/function
func IsFunction(data interface{}) bool {
	return typeIs(data, reflect.Func)
}

// IsInt will return true when type of the data is numeric integer
func IsInt(data interface{}) bool {
	return typeIs(data,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
	)
}

// IsMap will return true when type of the data is hash map
func IsMap(data interface{}) bool {
	return typeIs(data, reflect.Map)
}

// IsNil will return true when type of the data is nil
func IsNil(data interface{}) bool {
	if data == nil {
		return true
	}

	valueOfData := reflect.ValueOf(data)

	switch valueOfData.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer, reflect.Struct:
		if valueOfData.IsNil() {
			return true
		}
	}

	return false
}

// IsNumeric will return true when type of the data is numeric (float, uint, int)
func IsNumeric(data interface{}) bool {
	return typeIs(data,
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
		reflect.Uintptr,
	)
}

// IsPointer will return true when type of the data is pointer
func IsPointer(data interface{}) bool {
	return typeIs(data, reflect.Ptr)
}

// IsStructObject will return true when type of the data is object from struct
func IsStructObject(data interface{}) bool {
	return typeIs(data, reflect.Struct)
}

// IsTrue will return true when type of the data is bool, and the value is true
func IsTrue(data interface{}) bool {
	if data == nil {
		return false
	}

	if value, ok := data.(bool); ok {
		return value == true
	}

	return false
}

// IsUint will return true when type of the data is uint
func IsUint(data interface{}) bool {
	return typeIs(data,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
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
	if value, ok := data.(uintptr); ok {
		return value == 0
	}

	if value, ok := data.(complex64); ok {
		value128 := complex128(value)
		return math.Float64bits(real(value128)) == 0 && math.Float64bits(imag(value128)) == 0
	}
	if value, ok := data.(complex128); ok {
		return math.Float64bits(real(value)) == 0 && math.Float64bits(imag(value)) == 0
	}

	return false
}

// IsZeroValue reports whether value is the zero value for its type.
func IsZeroValue(data interface{}) bool {
	if data == nil {
		return true
	} else if value, ok := data.(string); ok {
		return value == ""
	} else if value, ok := data.(bool); ok {
		return value == false
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
	} else if value, ok := data.(uintptr); ok {
		return value == 0
	} else if value, ok := data.(complex64); ok {
		value128 := complex128(value)
		return math.Float64bits(real(value128)) == 0 && math.Float64bits(imag(value128)) == 0
	} else if value, ok := data.(complex128); ok {
		return math.Float64bits(real(value)) == 0 && math.Float64bits(imag(value)) == 0
	} else {
		if IsStructObject(data) {
			if IsNil(data) {
				return true
			}

			valueOfData := reflect.ValueOf(data)
			for i := 0; i < valueOfData.NumField(); i++ {
				if !IsZeroValue(valueOfData.Field(i).Interface()) {
					return false
				}
			}
		} else {
			if IsNil(data) {
				return true
			}
		}
	}

	return true
}

var IsEmpty = IsZeroValue
