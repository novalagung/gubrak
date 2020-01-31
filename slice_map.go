package gubrak

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Chunk function creates a slice of elements split into groups the length of `size`. If `data` can't be split evenly, the final chunk will be the remaining elements.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// 3 examples available:
func Chunk(data interface{}, size int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if !isZeroOrPositiveNumber(err, "size", size) {
			return nil
		}

		result := makeSlice(reflect.SliceOf(dataType))

		if dataValueLen == 0 {
			return result.Interface()
		}

		eachResult := makeSlice(dataType)

		if size > 0 {
			forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
				eachSize := eachResult.Len()
				if eachSize < size {
					eachResult = reflect.Append(eachResult, each)
				}

				if (eachSize+1) == size || (i+1) == dataValueLen {
					result = reflect.Append(result, eachResult)
					eachResult = reflect.MakeSlice(dataType, 0, 0)
				}
			})
		}

		return result.Interface()
	}(&err)

	return result, err
}

// Compact function creates a slice with all falsey values removed from the `data`. These values: `false`, `nil`, `0`, `""`, `(*string)(nil)`, and other nil-able types are considered to be falsey.
//
// Parameters
//
// This function requires one mandatory parameter:
//  data // type: slice, description: the slice to compact
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of filtered values
//  error // description: hold error message if there is an error
//
// Examples
//
// 4 examples available:
func Compact(data interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			target := each

			if target.Kind() == reflect.Ptr || target.Kind() == reflect.Interface {
				if target.IsNil() {
					return
				}

				target = target.Elem()
			}

			if target.Kind() == reflect.Ptr {
				if target.IsNil() {
					return
				}
			}

			ok := false

			switch target.Kind() {

			case reflect.Bool:
				ok = target.Bool()

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ok = target.Int() != 0

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				ok = target.Uint() != 0

			case reflect.Uintptr:
				ok = target.Elem().Uint() != 0

			case reflect.Float32, reflect.Float64:
				ok = target.Float() != 0

			case reflect.Complex64, reflect.Complex128:
				ok = target.Complex() != 0

			case reflect.Array:
				ok = target.Len() > 0

			case reflect.String:
				ok = target.String() != ""

			default:
				// case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
				ok = target.Interface() != nil
			}

			if ok {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Concat function creates a new slice concatenating `data` with any additional slices (the 2nd parameter and rest).
//
// Parameters
//
// This function requires one mandatory parameter `data`, and unlimited variadic parameters:
//  data        // type: slice, description: the slice to concatenate
//  dataConcat1 // type: slice, description: the values to concatenate
//  dataConcat2 // type: slice, description: the values to concatenate
//  dataConcat3 // type: slice, description: the values to concatenate
//  ...
//
// Return values
//
// This function return two values:
//  slice // description: returns the new concatenated slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Concat(data interface{}, dataConcats ...interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			result = reflect.Append(result, each)
		})

		for i, eachConcatenableData := range dataConcats {
			eachLabel := fmt.Sprintf("concat data %d", (i + 1))
			eachValue, eachType, _, eachValueLen := inspectData(eachConcatenableData)

			if !isSlice(err, eachLabel, eachValue) {
				continue
			}

			if dataValueLen == 0 {
				continue
			}

			if !isTypeEqual(err, "data", dataType, eachLabel, eachType) {
				continue
			}

			forEachSlice(eachValue, eachValueLen, func(each reflect.Value, i int) {
				result = reflect.Append(result, each)
			})
		}

		return result.Interface()
	}(&err)

	return result, err
}

// Count creates an object composed of keys generated from the results of running each element of `data` thru `iteratee`. The corresponding value of each key is the number of times the key was returned by `iteratee`.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data     // type: slice or map, description: the slice/map to iterate over
//  iteratee // optional, type: func(each anyType, i int)bool or func(value anyType, key anyType, i int), description: the function invoked per iteration.
//
// Return values
//
// This function return two values:
//  number // description: Returns the composed aggregate object
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Count(data interface{}, args ...interface{}) (int, error) {
	var err error

	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return 0
		}

		var callback interface{}
		if len(args) > 0 {
			callback = args[0]
		}

		dataValue, dataValueType, dataValueKind, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.Map {
				*err = nil
				return _countCollection(err, dataValue, dataValueType, dataValueKind, dataValueLen, callback)
			}

			return 0
		}

		return _countSlice(err, dataValue, dataValueType, dataValueKind, dataValueLen, callback)
	}(&err)

	return result, err
}

func _countSlice(err *error, dataValue reflect.Value, dataValueType reflect.Type, dataValueKind reflect.Kind, dataValueLen int, callback interface{}) int {

	var callbackValue reflect.Value
	var callbackType reflect.Type
	var callbackTypeNumIn int

	if callback != nil {
		callbackValue, callbackType = inspectFunc(err, callback)
		if *err != nil {
			return 0
		}

		callbackTypeNumIn = validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return 0
		}

		validateFuncOutputOneVarBool(err, callbackType, true)
		if *err != nil {
			return 0
		}
	}

	if dataValueLen == 0 {
		return 0
	}

	if callback == nil {
		return dataValueLen
	}

	resultCounter := 0

	forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
		res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
		if res[0].Bool() {
			resultCounter++
		}
	})

	return resultCounter
}

func _countCollection(err *error, dataValue reflect.Value, dataValueType reflect.Type, dataValueKind reflect.Kind, dataValueLen int, callback interface{}) int {

	var callbackValue reflect.Value
	var callbackType reflect.Type
	var callbackTypeNumIn int

	if callback != nil {
		callbackValue, callbackType = inspectFunc(err, callback)
		if *err != nil {
			return 0
		}

		callbackTypeNumIn = validateFuncInputForCollectionLoop(err, callbackType, dataValue)
		if *err != nil {
			return 0
		}

		validateFuncOutputOneVarBool(err, callbackType, true)
		if *err != nil {
			return 0
		}
	}

	if dataValueLen == 0 {
		return 0
	}

	if callback == nil {
		return dataValueLen
	}

	resultCounter := 0

	dataValueMapKeys := dataValue.MapKeys()
	forEachCollection(dataValue, dataValueMapKeys, func(value, key reflect.Value, i int) {
		res := callFuncCollectionLoop(callbackValue, value, key, callbackTypeNumIn)
		if res[0].Bool() {
			resultCounter++
		}
	})

	return resultCounter
}

// Difference function creates a slice of `data` values not included in the other given slices. The order and references of result values are determined by the first slice.
//
// Parameters
//
// This function requires one mandatory parameter `data`, and unlimited variadic parameters:
//  data   // type: slice, description: the slice to inspect
//  slice1 // optional, type: slice, description: the values to exclude
//  slice2 // optional, type: slice, description: the values to exclude
//  slice3 // optional, type: slice, description: the values to exclude
//  ...
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of filtered values
//  error // description: hold error message if there is an error
//
// Examples
//
// 3 examples available:
func Difference(data interface{}, compareData ...interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		compareDataMap := make(map[reflect.Value]int)
		for i, each := range compareData {
			eachLabel := fmt.Sprintf("difference data %d", (i + 1))
			eachValue, eachType, _, eachValueLen := inspectData(each)

			if !isSlice(err, eachLabel, eachValue) {
				continue
			}

			if dataValueLen == 0 {
				continue
			}

			if !isTypeEqual(err, "data", dataType, eachLabel, eachType) {
				continue
			}

			compareDataMap[eachValue] = eachValueLen
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			isFound := false

			for compareValue, compareValueLen := range compareDataMap {
				forEachSliceStoppable(compareValue, compareValueLen, func(inner reflect.Value, j int) bool {
					if each.Interface() == inner.Interface() {
						isFound = true
						return false
					}

					return true
				})
			}

			if !isFound {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Drop function creates a slice of `data` with `n` elements dropped from the beginning.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to query
//  n    // type: number, description: the number of elements to drop
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Drop(data interface{}, size int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if !isZeroOrPositiveNumber(err, "size", size) {
			return data
		}

		if size == 0 {
			return data
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			if i < size {
				return
			}

			result = reflect.Append(result, each)
		})

		return result.Interface()
	}(&err)

	return result, err
}

// DropRight function creates a slice of `data` with `n` elements dropped from the end.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to query
//  n    // type: number, description: the number of elements to drop
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func DropRight(data interface{}, size int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if !isZeroOrPositiveNumber(err, "size", size) {
			return data
		}

		if size == 0 {
			return data
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			if i < (dataValueLen - size) {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Each iterates over elements of `data` and invokes `iteratee` for each element. Iteratee functions may exit iteration early by explicitly returning false
//
// Parameters
//
// This function requires two mandatory parameters:
//  data     // type: slice or map, description: the slice/map to iterate over
//  iteratee // optional, type: FuncSliceLoopOutputBool, description: the function invoked per iteration. The second argument represents index of each element, and it's optional
//
// Return values
//
// This function return two values:
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Each(data, callback interface{}) error {
	var err error

	func(err *error) {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return
		}

		dataValue, dataValueType, dataValueKind, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.Map {
				*err = nil
				_eachCollection(err, dataValue, dataValueType, dataValueKind, dataValueLen, callback, true)
			}

			return
		}

		_eachSlice(err, dataValue, dataValueType, dataValueKind, dataValueLen, callback, true)
	}(&err)

	return err
}

func _eachSlice(err *error, dataValue reflect.Value, dataValueType reflect.Type, dataValueKind reflect.Kind, dataValueLen int, callback interface{}, isLoopIncremental bool) {
	callbackValue, callbackType := inspectFunc(err, callback)
	if *err != nil {
		return
	}

	callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
	if *err != nil {
		return
	}

	validateFuncOutputOneVarBool(err, callbackType, false)
	if *err != nil {
		return
	}

	if dataValueLen == 0 {
		return
	}

	forEachSliceStoppable(dataValue, dataValueLen, func(each reflect.Value, i int) bool {
		var res []reflect.Value
		if isLoopIncremental {
			res = callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
		} else {
			indexFromRight := dataValueLen - i - 1
			eachFromRight := dataValue.Index(indexFromRight)
			res = callFuncSliceLoop(callbackValue, eachFromRight, indexFromRight, callbackTypeNumIn)
		}

		if len(res) > 0 {
			return res[0].Bool()
		}

		return true
	})
}

func _eachCollection(err *error, dataValue reflect.Value, dataValueType reflect.Type, dataValueKind reflect.Kind, dataValueLen int, callback interface{}, isLoopIncremental bool) {
	callbackValue, callbackType := inspectFunc(err, callback)
	if *err != nil {
		return
	}

	callbackTypeNumIn := validateFuncInputForCollectionLoop(err, callbackType, dataValue)
	if *err != nil {
		return
	}

	validateFuncOutputNone(err, callbackType)
	if *err != nil {
		return
	}

	if dataValueLen == 0 {
		return
	}

	dataValueMapKeys := dataValue.MapKeys()
	forEachCollectionStoppable(dataValue, dataValueMapKeys, func(value reflect.Value, key reflect.Value, i int) bool {
		var res []reflect.Value
		if isLoopIncremental {
			res = callFuncCollectionLoop(callbackValue, value, key, callbackTypeNumIn)
		} else {
			indexFromRight := dataValueLen - i - 1
			keyFromRight := dataValueMapKeys[indexFromRight]
			valueFromRight := dataValue.MapIndex(keyFromRight)
			res = callFuncCollectionLoop(callbackValue, valueFromRight, keyFromRight, callbackTypeNumIn)
		}

		if len(res) > 0 {
			return res[0].Bool()
		}

		return true
	})
}

// EachRight function is like ForEach() except that it iterates over elements of collection from right to left.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func EachRight(data, callback interface{}) error {
	var err error

	func(err *error) {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return
		}

		dataValue, dataValueType, dataValueKind, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.Map {
				*err = nil
				_eachCollection(err, dataValue, dataValueType, dataValueKind, dataValueLen, callback, false)
			}

			return
		}

		_eachSlice(err, dataValue, dataValueType, dataValueKind, dataValueLen, callback, false)
	}(&err)

	return err
}

// Fill function fills elements of `data` with `value` from `start` up to, but not including, `end`.
//
// Parameters
//
// This function requires two mandatory parameters `data` and `value`; and two optional parameters:
//  data          // type: slice, description: the slice to fill
//  value         // type: anyType, description: the value to fill slice with. This variable's data type must be same with slice's element data type
//  start=0       // optional, type: number, description: the start position
//  end=len(data) // optional, type: number, description: the end position
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 3 examples available:
func Fill(data, value interface{}, args ...int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		startIndex := 0
		lastIndex := dataValueLen

		if len(args) > 0 {
			startIndex = args[0]

			if len(args) > 1 {
				lastIndex = args[1]
			}
		}

		if !isZeroOrPositiveNumber(err, "start index", startIndex) {
			return data
		}

		if !isZeroOrPositiveNumber(err, "last index", lastIndex) {
			return data
		}

		if !isLeftShouldBeGreaterOrEqualThanRight(err, "last index", lastIndex, "start index", startIndex) {
			return nil
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		fillValue, fillType, _, _ := inspectData(value)

		if fillType.Kind() != dataType.Elem().Kind() {
			*err = errors.New("replacement data type must be same with slice's element type")
			return nil
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			if i >= startIndex && i < lastIndex {
				result = reflect.Append(result, fillValue)
			} else {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Filter function iterates over elements of collection, returning an array of all elements predicate returns truthy for. The predicate is invoked with two arguments: (value, index).
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Filter(data, callback interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, dataValueKind, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.Map {
				*err = nil
				return _filterCollection(err, dataValue, dataType, dataValueKind, dataValueLen, callback)
			}

			return nil
		}

		return _filterSlice(err, dataValue, dataType, dataValueKind, dataValueLen, callback)
	}(&err)

	return result, err
}

func _filterSlice(err *error, dataValue reflect.Value, dataValueType reflect.Type, dataValueKind reflect.Kind, dataValueLen int, callback interface{}) interface{} {
	callbackValue, callbackType := inspectFunc(err, callback)
	if *err != nil {
		return nil
	}

	callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
	if *err != nil {
		return nil
	}

	validateFuncOutputOneVarBool(err, callbackType, true)
	if *err != nil {
		return nil
	}

	result := makeSlice(dataValueType)

	if dataValueLen == 0 {
		return result.Interface()
	}

	forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
		res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
		if res[0].Bool() {
			result = reflect.Append(result, each)
		}
	})

	return result.Interface()
}

func _filterCollection(err *error, dataValue reflect.Value, dataValueType reflect.Type, dataValueKind reflect.Kind, dataValueLen int, callback interface{}) interface{} {
	callbackValue, callbackType := inspectFunc(err, callback)
	if *err != nil {
		return nil
	}

	callbackTypeNumIn := validateFuncInputForCollectionLoop(err, callbackType, dataValue)
	if *err != nil {
		return nil
	}

	validateFuncOutputOneVarBool(err, callbackType, true)
	if *err != nil {
		return nil
	}

	result := reflect.MakeMap(reflect.MapOf(dataValueType.Key(), dataValueType.Elem()))

	if dataValueLen == 0 {
		return result.Interface()
	}

	dataValueMapKeys := dataValue.MapKeys()
	forEachCollection(dataValue, dataValueMapKeys, func(value reflect.Value, key reflect.Value, i int) {
		res := callFuncCollectionLoop(callbackValue, value, key, callbackTypeNumIn)
		if res[0].Bool() {
			result.SetMapIndex(key, value)
		}
	})

	return result.Interface()
}

// Find function iterates over elements of collection, returning the first element predicate returns truthy for. The predicate is invoked with three arguments: (value, index).
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Find(data, callback interface{}, args ...int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, callback)
		if *err != nil {
			return nil
		}

		callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return nil
		}

		validateFuncOutputOneVarBool(err, callbackType, true)
		if *err != nil {
			return nil
		}

		fromIndex := 0
		if len(args) > 0 {
			fromIndex = args[0]
		}

		if !isZeroOrPositiveNumber(err, "from index", fromIndex) {
			return nil
		}

		if dataValueLen == 0 {
			return nil
		}

		isFound := false
		result := reflect.New(dataType)

		forEachSliceStoppable(dataValue, dataValueLen, func(each reflect.Value, i int) bool {
			if i < fromIndex {
				return true
			}

			res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
			if res[0].Bool() {
				isFound = true
				result = each
				return false
			}

			return true
		})

		if isFound {
			return result.Interface()
		}

		return nil
	}(&err)

	return result, err
}

// FindIndex function is similar like `Find`, except that it returns the index of the first element `predicate` returns truthy for, instead of the element itself.
//
// Parameters
//
// This function requires two mandatory parameters `data` and `predicate`; and two optional parameters:
//  data        // type: slice, description: the slice to inspect
//  predicate   // type: func(each anyType, i int)bool, description: the function invoked per iteration.
//  fromIndex=0 // optional, type: number, description: the index to search from
//
// Return values
//
// This function return two values:
//  number // description: returns the index of the found element, else `-1`
//  error  // description: hold error message if there is an error
//
// Examples
//
// 5 examples available:
func FindIndex(data, predicate interface{}, args ...int) (int, error) {
	var err error

	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return -1
		}

		dataValue, _, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return -1
		}

		callbackValue, callbackType := inspectFunc(err, predicate)
		if *err != nil {
			return -1
		}

		callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return -1
		}

		validateFuncOutputOneVarBool(err, callbackType, true)
		if *err != nil {
			return -1
		}

		if dataValueLen == 0 {
			return -1
		}

		startIndex := 0
		if len(args) > 0 {
			startIndex = args[0]
		}

		result := -1

		forEachSliceStoppable(dataValue, dataValueLen, func(each reflect.Value, i int) bool {
			if i < startIndex {
				return true
			}

			res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
			if res[0].Bool() {
				result = i
				return false
			}

			return true
		})

		return result
	}(&err)

	return result, err
}

// FindLast function is like Find() except that it iterates over elements of collection from right to left.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func FindLast(data, callback interface{}, args ...int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, callback)
		if *err != nil {
			return nil
		}

		callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return nil
		}

		validateFuncOutputOneVarBool(err, callbackType, true)
		if *err != nil {
			return nil
		}

		lastIndex := dataValueLen - 1
		if len(args) > 0 {
			lastIndex = args[0]
		}

		if !isZeroOrPositiveNumber(err, "last index", lastIndex) {
			return nil
		}

		if dataValueLen == 0 {
			return nil
		}

		isFound := false
		result := reflect.New(dataType)

		forEachSliceStoppable(dataValue, dataValueLen, func(each reflect.Value, i int) bool {
			reverseIndex := dataValueLen - 1 - i
			if reverseIndex > lastIndex {
				return true
			}

			eachReverse := dataValue.Index(reverseIndex)
			res := callFuncSliceLoop(callbackValue, dataValue.Index(reverseIndex), reverseIndex, callbackTypeNumIn)
			if res[0].Bool() {
				isFound = true
				result = eachReverse
				return false
			}

			return true
		})

		if isFound {
			return result.Interface()
		}

		return nil
	}(&err)

	return result, err
}

// FindLastIndex function is similar like `FindIndex`, except that it iterates over elements of `data` from right to left.
//
// Parameters
//
// This function requires two mandatory parameters `data` and `predicate`; and an optional parameter:
//  data                  // type: slice, description: the slice to inspect
//  predicate             // type: func(each anyType, i int)bool, description: the function invoked per iteration.it's optional
//  fromIndex=len(data)-1 // optional, type: number, description: the index to search from
//
// Return values
//
// This function return two values:
//  number // description: returns the index of the found element, else `-1`
//  error  // description: hold error message if there is an error
//
// Examples
//
// 4 examples available:
func FindLastIndex(data, predicate interface{}, args ...int) (int, error) {
	var err error

	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return -1
		}

		dataValue, _, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return -1
		}

		callbackValue, callbackType := inspectFunc(err, predicate)
		if *err != nil {
			return -1
		}

		callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return -1
		}

		validateFuncOutputOneVarBool(err, callbackType, true)
		if *err != nil {
			return -1
		}

		if dataValueLen == 0 {
			return -1
		}

		endIndex := dataValueLen
		if len(args) > 0 {
			endIndex = args[0]
		}

		result := -1

		forEachSliceStoppable(dataValue, dataValueLen, func(each reflect.Value, i int) bool {
			if i > endIndex {
				return true
			}

			res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
			if res[0].Bool() {
				result = i
			}

			return true
		})

		if result < 0 {
			return result
		}

		return result
	}(&err)

	return result, err
}

// First function gets the first element of `data`.
//
// Parameters
//
// This function requires one mandatory parameter `data`:
//  data // type: slice, description: the slice to query
//
// Return values
//
// This function return two values:
//  anyType // description: returns the first element of data
//  error   // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func First(data interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, _, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return nil
		}

		return dataValue.Index(0).Interface()
	}(&err)

	return result, err
}

// ForEach is alias of Each()
func ForEach(data, callback interface{}) error {
	return Each(data, callback)
}

// ForEachRight is alias of EachRight()
func ForEachRight(data, callback interface{}) error {
	return EachRight(data, callback)
}

// FromPairs function returns an object composed from key-value `data`.
//
// Parameters
//
// This function requires one mandatory parameter `data`:
//  data // type: [][]interface{}, description: the key-value pairs
//
// Return values
//
// This function return two values:
//  map[interface{}]interface{} // description: returns the new object
//  error                       // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func FromPairs(data interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataValueType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueType.Elem().Kind() != reflect.Interface {
			*err = errors.New("supported type only []interface{}")
			return nil
		}

		result := make(map[interface{}]interface{}, 0)

		if dataValueLen == 0 {
			return result
		}

		forEachSliceStoppable(dataValue, dataValueLen, func(each reflect.Value, i int) bool {
			if *err != nil {
				result = nil
				return false
			}

			eachSlice := each.Elem()
			eachSliceLen := eachSlice.Len()

			if eachSliceLen > 2 {
				eachSliceLen = 2
			}

			if eachSliceLen > 0 {
				eachSliceKey := eachSlice.Index(0).Interface()
				result[eachSliceKey] = nil

				if eachSliceLen > 1 {
					eachSliceVal := eachSlice.Index(1).Interface()
					result[eachSliceKey] = eachSliceVal
				}
			}

			return true
		})

		return result
	}(&err)

	return result, err
}

// GroupBy function creates an object composed of keys generated from the results of running each element of collection thru iteratee. The order of grouped values is determined by the order they occur in collection. The corresponding value of each key is an array of elements responsible for generating the key. The iteratee is invoked with one argument: (value).
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func GroupBy(data, callback interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, callback)
		if *err != nil {
			return nil
		}

		callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return nil
		}

		validateFuncOutputOneVarDynamic(err, callbackType)
		if *err != nil {
			return nil
		}

		result := reflect.MakeMap(reflect.MapOf(callbackType.Out(0), dataType))

		if dataValueLen == 0 {
			return result.Interface()
		}

		resultMap := make(map[interface{}]reflect.Value)

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
			resActualValue := res[0].Interface()

			if _, ok := resultMap[resActualValue]; !ok {
				resultMap[resActualValue] = makeSlice(dataType)
			}

			mapValues := resultMap[resActualValue]
			mapValues = reflect.Append(mapValues, each)
			resultMap[resActualValue] = mapValues
			result.SetMapIndex(res[0], mapValues)
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Head function is an alias of `First`.
func Head(data interface{}) (interface{}, error) {
	return First(data)
}

// Includes function checks if value is in collection. If collection is a string, it's checked for a substring of value, otherwise SameValueZero is used for equality comparisons. If fromIndex is negative, it's used as the offset from the end of collection.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Includes(data, search interface{}, args ...int) (bool, error) {
	var err error

	result := func(err *error) bool {
		defer catch(err)

		if dataValue, dataOK := data.(string); dataOK {
			if searchValue, searchOK := search.(string); searchOK {
				return strings.Contains(dataValue, searchValue)
			}

			return false
		}

		if !isNonNilData(err, "data", data) {
			return false
		}

		dataValue, _, dataValueKind, dataValueLen := inspectData(data)

		startIndex := 0
		if len(args) > 0 {
			startIndex = args[0]
		}

		if !isZeroOrPositiveNumber(err, "start index", startIndex) {
			return false
		}

		if dataValueLen == 0 {
			return false
		}

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.Map {
				*err = nil
				return _includesCollection(err, dataValue, search, startIndex)
			}

			*err = errors.New((*err).Error() + ", map, or a string")
			return false
		}

		return _includesSlice(err, dataValue, dataValueLen, search, startIndex)
	}(&err)

	return result, err
}

func _includesSlice(err *error, dataValue reflect.Value, dataValueLen int, search interface{}, startIndex int) bool {
	isFound := false

	forEachSliceStoppable(dataValue, dataValueLen, func(each reflect.Value, i int) bool {
		if i < startIndex {
			return true
		}

		eachActualValue := each.Interface()

		if eachActualValue == search {
			isFound = true
			return false
		}

		return true
	})

	return isFound
}

func _includesCollection(err *error, dataValue reflect.Value, search interface{}, startIndex int) bool {
	isFound := false
	counter := 0

	dataValueMapKeys := dataValue.MapKeys()
	forEachCollectionStoppable(dataValue, dataValueMapKeys, func(value reflect.Value, key reflect.Value, i int) bool {
		defer func() {
			counter++
		}()

		if counter < startIndex {
			return true
		}

		eachActualValue := value.Interface()

		if eachActualValue == search && !isFound {
			isFound = true
			return false
		}

		return true
	})

	return isFound
}

// IndexOf function gets the index at which the first occurrence of `search` is found in `data`. If `fromIndex` is negative, it's used as the offset from the end of `data`.
//
// Parameters
//
// This function requires two mandatory parameters `data` and `value`; and one optional parameter:
//  data        // type: slice, description: the slice to inspect
//  value       // type: anyType, description: the value to search for
//  fromIndex=0 // optional, type: number, description: the index to search from
//
// Return values
//
// This function return two values:
//  number // description: returns the index of the matched value, else -1
//  error  // description: hold error message if there is an error
//
// Examples
//
// 3 examples available:
func IndexOf(data interface{}, search interface{}, args ...int) (int, error) {
	var err error

	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return -1
		}

		dataValue, _, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return -1
		}

		startIndex := 0
		if len(args) > 0 {
			startIndex = args[0]
		}

		if startIndex >= dataValueLen {
			return -1
		}

		if dataValueLen == 0 {
			return -1
		}

		result := -1

		forEachSliceStoppable(dataValue, dataValueLen, func(each reflect.Value, i int) bool {
			if startIndex > -1 {
				if startIndex > 0 && i < startIndex {
					return true
				}

				if each.Interface() == search && result == -1 {
					result = i
					return false
				}
			} else {
				if i > (startIndex*-1)-1 {
					return true
				}

				iFromRight := dataValueLen - i - 1
				eachFromRight := dataValue.Index(iFromRight)

				if eachFromRight.Interface() == search {
					result = iFromRight
					return true
				}
			}

			return true
		})

		return result
	}(&err)

	return result, err
}

// Initial function gets all but the last element of `data`.
//
// Parameters
//
// This function requires one mandatory parameter:
//  data // type: slice, description: the slice to query
//
// Return values
//
// This function return two values:
//  slice // description: returns the slice of `data`
//  error // description: hold error message if there is an error
//
// Examples
//
// 4 examples available:
func Initial(data interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return makeSlice(dataType).Interface()
		}

		return dataValue.Slice(0, dataValueLen-1).Interface()
	}(&err)

	return result, err
}

// Intersection function creates a slice of unique values that are included in all given slices. The order and references of result values are determined by the first slice.
//
// Parameters
//
// This function requires one mandatory parameter `data`; and unlimited variadic parameters:
//  data           // type: slice, the slice to inspect
//  dataIntersect1 // optional, type: slice, the values to compare
//  dataIntersect2 // optional, type: slice, the values to compare
//  dataIntersect3 // optional, type: slice, the values to compare
// ...
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of intersecting values
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Intersection(data interface{}, dataIntersects ...interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		type CompareMap struct {
			Value reflect.Value
			Len   int
		}
		compareValueInReflect := make([]CompareMap, 0)

		for _, compare := range dataIntersects {
			eachValue, _, _, eachValueLen := inspectData(compare)

			if isSlice(err, "data", eachValue) {
				compareValueInReflect = append(compareValueInReflect, CompareMap{
					Value: eachValue,
					Len:   eachValueLen,
				})
			} else {
				*err = errors.New("All data should be slice")
				return nil
			}
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		resultMap := make(map[interface{}]bool)

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			eachActualValue := each.Interface()

			isValueExists := true
			isJustStarted := true

			for _, eachCompare := range compareValueInReflect {
				isInnerExists := false

				forEachSliceStoppable(eachCompare.Value, eachCompare.Len, func(inner reflect.Value, j int) bool {
					if eachActualValue == inner.Interface() {
						isInnerExists = true
						return false
					}

					return true
				})

				if isJustStarted {
					isValueExists = isInnerExists
				} else {
					isValueExists = isValueExists && isInnerExists
				}
			}

			if isValueExists {
				if _, ok := resultMap[eachActualValue]; !ok {
					resultMap[eachActualValue] = true
					result = reflect.Append(result, each)
				}
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Join function converts all elements in `data` into a string separated by `separator`.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data      // type: slice, description: the slices to convert
//  separator // type: string, description: the element separator
//
// Return values
//
// This function return two values:
//  string // description: returns the joined string
//  error  // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Join(data interface{}, separator string) (string, error) {
	var err error

	result := func(err *error) string {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return ""
		}

		dataValue, _, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return ""
		}

		if val, ok := data.([]string); ok {
			return strings.Join(val, separator)
		}

		if dataValueLen == 0 {
			return ""
		}

		dataInStringSlice := make([]string, 0)

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {

			if each.Interface() != nil {
				target := each

				if each.Kind() == reflect.Interface {
					target = each.Elem()
				}

				if target.Kind() == reflect.String {
					dataInStringSlice = append(dataInStringSlice, target.String())
				} else {
					dataInStringSlice = append(dataInStringSlice, fmt.Sprintf("%v", target.Interface()))
				}
			}
		})

		if len(dataInStringSlice) == 0 {
			return ""
		}

		return strings.Join(dataInStringSlice, separator)
	}(&err)

	return result, err
}

// KeyBy function creates an object composed of keys generated from the results of running each element of collection thru iteratee. The corresponding value of each key is the last element responsible for generating the key. The iteratee is invoked with one argument: (value).
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func KeyBy(data, callback interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataValueType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, callback)
		if *err != nil {
			return nil
		}

		callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return nil
		}

		validateFuncOutputOneVarDynamic(err, callbackType)
		if *err != nil {
			return nil
		}

		valueElemType := dataValueType.Elem()
		result := reflect.MakeMap(reflect.MapOf(callbackType.Out(0), valueElemType))

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
			result.SetMapIndex(res[0], each)
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Last function gets the last element of `data`.
//
// Parameters
//
// This function requires one mandatory parameter:
//  data // type: slice, description: the slices to query
//
// Return values
//
// This function return two values:
//  anyType // description: returns the last element of `data`
//  error   // description: hold error message if there is an error
//
// Examples
//
// 3 examples available:
func Last(data interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, _, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return nil
		}

		return dataValue.Index(dataValueLen - 1).Interface()
	}(&err)

	return result, err
}

// LastIndexOf function is like `IndexOf`, except that it iterates over elements of `data` from right to left.
//
// Parameters
//
// This function requires two mandatory parameters `data` and `search`; and an optional parameter:
//  data                  // type: slice, description: the slices to inspect
//  search                // type: anyType, description: the value to search for
//  fromIndex=len(data)-1 // type: number, description: the index to search from
//
// Return values
//
// This function return two values:
//  number // description: returns the index of the matched value, else `-1`
//  error  // description: hold error message if there is an error
//
// Examples
//
// 3 examples available:
func LastIndexOf(data interface{}, search interface{}, args ...int) (int, error) {
	var err error

	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return -1
		}

		dataValue, _, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return -1
		}

		startIndex := dataValueLen - 1
		if len(args) > 0 {
			startIndex = args[0]
		}

		if dataValueLen == 0 {
			return -1
		}

		result := -1

		forEachSliceStoppable(dataValue, dataValueLen, func(each reflect.Value, i int) bool {
			if startIndex > -1 {
				iFromRight := startIndex - i
				if iFromRight > (dataValueLen-1) || iFromRight < 0 {
					return true
				}

				eachFromRight := dataValue.Index(iFromRight)
				if eachFromRight.Interface() == search && result == -1 {
					result = iFromRight
					return true
				}
			} else {
				iFromRight := dataValueLen + startIndex - i
				if iFromRight > (dataValueLen-1) || iFromRight < 0 {
					return true
				}

				eachFromRight := dataValue.Index(iFromRight)
				if eachFromRight.Interface() == search && result == -1 {
					result = iFromRight
					return true
				}
			}

			return true
		})

		return result
	}(&err)

	return result, err
}

// Map function creates an array of values by running each element in collection thru iteratee. The iteratee is invoked with two arguments: (value, index).
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Map(data, callback interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, _, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, callback)
		if *err != nil {
			return nil
		}

		callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return nil
		}

		validateFuncOutputOneVarDynamic(err, callbackType)
		if *err != nil {
			return nil
		}

		result := makeSlice(reflect.SliceOf(callbackType.Out(0)))

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
			result = reflect.Append(result, res[0])
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Nth function gets the element at index `n` of `data`. If `n` is negative, the nth element from the end is returned.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slices to query
//  n=0  // type: number, description: The index of the element to return
//
// Return values
//
// This function return two values:
//  AnyType // description: returns the nth element of `data`
//  error   // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Nth(data interface{}, i int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, _, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return nil
		}

		if i < 0 {
			i = dataValueLen + i
		}

		if i < dataValueLen {
			return dataValue.Index(i).Interface()
		}

		return nil
	}(&err)

	return result, err
}

// OrderBy sort slices. If orders is unspecified, all values are sorted in ascending order. Otherwise, specify an order of "desc" for descending or "asc" for ascending sort order of corresponding values. The algorithm used is merge sort, as per savigo's post on https://sagivo.com/go-sort-faster-4869bdabc670
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func OrderBy(data, callback interface{}, args ...bool) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataValueType, _, _ := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, callback)
		if *err != nil {
			return nil
		}

		validateFuncInputForSliceLoopWithoutIndex(err, callbackType, dataValue)
		if *err != nil {
			return nil
		}

		validateFuncOutputOneVarDynamic(err, callbackType)
		if *err != nil {
			return nil
		}

		isAscending := true
		if len(args) > 0 {
			isAscending = args[0]
		}

		isAsync := false
		if len(args) > 1 {
			isAsync = args[1]
		}

		// =====

		var _doSortAsync func(reflect.Value, chan reflect.Value)
		var _doSortSync func(reflect.Value) reflect.Value
		var _doMerge func(reflect.Value, reflect.Value) reflect.Value

		_doSortAsync = func(slice reflect.Value, c chan reflect.Value) {
			sliceLen := slice.Len()

			if sliceLen < -1 {
				c <- _doSortSync(slice)
				return
			}

			if sliceLen < 2 {
				c <- slice
				return
			}

			mid := sliceLen / 2

			c1 := make(chan reflect.Value, 1)
			c2 := make(chan reflect.Value, 1)

			go _doSortAsync(slice.Slice(0, mid), c1)
			go _doSortAsync(slice.Slice(mid, sliceLen), c2)

			go func() {
				c <- _doMerge(<-c1, <-c2)
			}()
		}

		_doSortSync = func(slice reflect.Value) reflect.Value {
			sliceLen := slice.Len()
			if sliceLen < 2 {
				return slice
			}

			mid := sliceLen / 2
			leftSlice := _doSortSync(slice.Slice(0, mid))
			rightSlice := _doSortSync(slice.Slice(mid, slice.Len()))

			return _doMerge(leftSlice, rightSlice)
		}

		_doMerge = func(leftSlice, rightSlice reflect.Value) reflect.Value {
			bitSize := 64
			base := 10

			resultLen := leftSlice.Len() + rightSlice.Len()
			result := makeSlice(dataValueType, resultLen, resultLen)

			isSortable := true

			var i, j int

			for isSortable && i < leftSlice.Len() && j < rightSlice.Len() {
				isLeftLowerThanRight := false

				leftElem := leftSlice.Index(i)
				leftValue := callFuncSliceLoop(callbackValue, leftElem, i, 1)[0]

				rightElem := rightSlice.Index(j)
				rightValue := callFuncSliceLoop(callbackValue, rightElem, j, 1)[0]

				switch leftValue.Kind() {
				case reflect.String:
					if rightValue.Kind() == reflect.String {
						isLeftLowerThanRight = leftValue.String() <= rightValue.String()
					} else {
						isLeftLowerThanRight = leftValue.String() <= fmt.Sprintf("%v", rightValue.Interface())
					}

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

					switch rightValue.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						isLeftLowerThanRight = leftValue.Int() <= rightValue.Int()

					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						isLeftLowerThanRight = uint64(leftValue.Int()) <= rightValue.Uint()

					case reflect.Float32, reflect.Float64:
						isLeftLowerThanRight = float64(leftValue.Int()) <= rightValue.Float()

					case reflect.String:
						v, _ := strconv.ParseInt(fmt.Sprintf("%v", rightValue.Interface()), base, bitSize)
						isLeftLowerThanRight = leftValue.Int() <= v

					default:
						isSortable = false
					}

				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

					switch rightValue.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						isLeftLowerThanRight = leftValue.Uint() <= uint64(rightValue.Int())

					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						isLeftLowerThanRight = leftValue.Uint() <= rightValue.Uint()

					case reflect.Float32, reflect.Float64:
						isLeftLowerThanRight = float64(leftValue.Uint()) <= rightValue.Float()

					case reflect.String:
						v, _ := strconv.ParseUint(fmt.Sprintf("%v", rightValue.Interface()), base, bitSize)
						isLeftLowerThanRight = leftValue.Uint() <= v

					default:
						isSortable = false
					}

				case reflect.Float32, reflect.Float64:

					switch rightValue.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						isLeftLowerThanRight = leftValue.Float() <= float64(rightValue.Int())

					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						isLeftLowerThanRight = leftValue.Float() <= float64(rightValue.Uint())

					case reflect.Float32, reflect.Float64:
						isLeftLowerThanRight = leftValue.Float() <= rightValue.Float()

					case reflect.String:
						s := strings.TrimSpace(fmt.Sprintf("%v", rightValue.Interface()))
						v := float64(0)
						if s != "" {
							var errConvertion error
							v, errConvertion = strconv.ParseFloat(s, bitSize)
							if errConvertion != nil {
								v = 0
							}
						}
						isLeftLowerThanRight = leftValue.Float() <= v

					default:
						isSortable = false
					}

				default:
					isSortable = false
				}

				if isLeftLowerThanRight {
					if isAscending {
						result.Index(i + j).Set(leftElem)
						i++
					} else {
						result.Index(i + j).Set(rightElem)
						j++
					}
				} else {
					if isAscending {
						result.Index(i + j).Set(rightElem)
						j++
					} else {
						result.Index(i + j).Set(leftElem)
						i++
					}
				}
			}

			if !isSortable {
				return dataValue
			}

			for i < leftSlice.Len() {
				result.Index(i + j).Set(leftSlice.Index(i))
				i++
			}

			for j < rightSlice.Len() {
				result.Index(i + j).Set(rightSlice.Index(j))
				j++
			}

			return result
		}

		if isAsync {
			c := make(chan reflect.Value)
			_doSortAsync(dataValue, c)

			return (<-c).Interface()
		}

		return _doSortSync(dataValue).Interface()
	}(&err)

	return result, err
}

// Partition function creates an array of elements split into two groups, the first of which contains elements predicate returns truthy for, the second of which contains elements predicate returns falsey for. The predicate is invoked with one argument: (value).
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Partition(data, callback interface{}) (interface{}, interface{}, error) {
	var truhty, falsey interface{}
	var err error

	truhty, falsey = func(err *error) (interface{}, interface{}) {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil, nil
		}

		dataValue, dataValueType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil, nil
		}

		callbackValue, callbackType := inspectFunc(err, callback)
		if *err != nil {
			return nil, nil
		}

		callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return nil, nil
		}

		validateFuncOutputOneVarBool(err, callbackType, true)
		if *err != nil {
			return nil, nil
		}

		resultTruhty := makeSlice(dataValueType)
		resultFalsey := makeSlice(dataValueType)

		if dataValueLen == 0 {
			return resultTruhty, resultFalsey
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)

			if res[0].Bool() {
				resultTruhty = reflect.Append(resultTruhty, each)
			} else {
				resultFalsey = reflect.Append(resultFalsey, each)
			}
		})

		return resultTruhty.Interface(), resultFalsey.Interface()
	}(&err)

	return truhty, falsey, err
}

// Pull function removes all given values from `data`.
//
// Parameters
//
// This function requires one mandatory parameter `data`; and unlimited variadic parameters:
//  data  // type: slice, description: the slices to modify
//  item1 // optional, type: anyType, description: item to remove
//  item2 // optional, type: anyType, description: item to remove
//  item3 // optional, type: anyType, description: item to remove
//  ...
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 3 examples available:
func Pull(data interface{}, items ...interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, valueType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if len(items) == 0 {
			return data
		}

		result := makeSlice(valueType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			eachRealValue := each.Interface()
			isFound := false

			for _, item := range items {
				if item == eachRealValue {
					isFound = true
				}
			}

			if !isFound {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// PullAll function is similar like `Pull`, except that it accepts a slice of values to remove.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data  // type: slice, description: the slices to modify
//  items // type: slice, description: items to remove
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func PullAll(data interface{}, items interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, valueType, _, dataValueLen := inspectData(data)
		itemsValue, _, _, itemsValueLen := inspectData(items)

		if !isSlice(err, "data", dataValue, itemsValue) {
			return nil
		}

		if itemsValueLen == 0 {
			return data
		}

		result := makeSlice(valueType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			eachRealValue := each.Interface()
			isFound := false

			forEachSliceStoppable(itemsValue, itemsValueLen, func(item reflect.Value, i int) bool {
				itemRealValue := item.Interface()

				if itemRealValue == eachRealValue {
					isFound = true
					return false
				}

				return true
			})

			if !isFound {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// PullAt function removes elements from `data` corresponding to `indexes` and returns an array of removed elements.
//
// Parameters
//
// This function requires one mandatory parameter `data`; and unlimited variadic parameters:
//  data     // type: slice, description: the slices to modify
//  indexes1 // optional, type: int, description: index to remove
//  indexes2 // optional, type: int, description: index to remove
//  indexes3 // optional, type: int, description: index to remove
//  ...
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 1 examples available:
func PullAt(data interface{}, indexes ...int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		for _, index := range indexes {
			if !isZeroOrPositiveNumber(err, "index", index) {
				return data
			}
		}

		if len(indexes) == 0 {
			return data
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			isFound := false

			for _, index := range indexes {
				if index == i {
					isFound = true
				}
			}

			if !isFound {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Reduce function reduces collection to a value which is the accumulated result of running each element in collection thru iteratee, where each successive invocation is supplied the return value of the previous. If accumulator is not given, the first element of collection is used as the initial value. The iteratee is invoked with four arguments: (accumulator, value, index|key, collection)
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Reduce(data, callback, initial interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataValueType, dataValueKind, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.Map {
				*err = nil
				return _reduceCollection(err, dataValue, dataValueType, dataValueKind, dataValueLen, callback, initial)
			}

			return nil
		}

		return _reduceSlice(err, dataValue, dataValueType, dataValueKind, dataValueLen, callback, initial)
	}(&err)

	return result, err
}

func _reduceCollection(err *error, dataValue reflect.Value, dataValueType reflect.Type, dataValueKind reflect.Kind, dataValueLen int, callback, initial interface{}) interface{} {

	callbackValue, callbackType := inspectFunc(err, callback)
	if *err != nil {
		return nil
	}

	initialValue, initialValueType, _, _ := inspectData(initial)

	callbackValueNumIn := callbackType.NumIn()
	if callbackValueNumIn < 2 || callbackValueNumIn > 3 {
		*err = errors.New("callback must only have two or three parameters")
		return nil
	}

	if callbackType.In(0).Kind() != initialValueType.Kind() {
		*err = errors.New("callback 1st parameter's data type should be same with initial value's data type")
		return nil
	}

	if callbackType.In(1).Kind() != dataValueType.Elem().Kind() {
		*err = errors.New("callback 2nd parameter's data type should be same with map value data type")
		return nil
	}

	if callbackValueNumIn > 2 {
		if callbackType.In(2).Kind() != dataValueType.Key().Kind() {
			*err = errors.New("callback 3rd parameter's data type should be same with map key type")
			return nil
		}
	}

	validateFuncOutputOneVarDynamic(err, callbackType)
	if *err != nil {
		return nil
	}

	result := initialValue

	dataValueMapKeys := dataValue.MapKeys()
	forEachCollection(dataValue, dataValueMapKeys, func(value, key reflect.Value, i int) {
		if callbackValueNumIn == 2 {
			result = callbackValue.Call([]reflect.Value{result, value})[0]
		} else {
			result = callbackValue.Call([]reflect.Value{result, value, key})[0]
		}
	})

	return result.Interface()
}

func _reduceSlice(err *error, dataValue reflect.Value, dataValueType reflect.Type, dataValueKind reflect.Kind, dataValueLen int, callback, initial interface{}) interface{} {

	callbackValue, callbackType := inspectFunc(err, callback)
	if *err != nil {
		return nil
	}

	initialValue, initialValueType, _, _ := inspectData(initial)

	callbackValueNumIn := callbackType.NumIn()
	if callbackValueNumIn < 2 || callbackValueNumIn > 3 {
		*err = errors.New("callback must only have two or three parameters")
		return nil
	}

	if callbackType.In(0).Kind() != initialValueType.Kind() {
		*err = errors.New("callback 1st parameter's data type should be same with initial value's data type")
		return nil
	}

	if callbackType.In(1).Kind() != dataValue.Index(0).Kind() {
		*err = errors.New("callback 2nd parameter's data type should be same with slice element data type")
		return nil
	}

	if callbackValueNumIn > 2 {
		if callbackType.In(2).Kind() != reflect.Int {
			*err = errors.New("callback 3rd parameter's data type should be int")
			return nil
		}
	}

	validateFuncOutputOneVarDynamic(err, callbackType)
	if *err != nil {
		return nil
	}

	result := initialValue

	forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
		if callbackValueNumIn == 2 {
			result = callbackValue.Call([]reflect.Value{result, each})[0]
		} else {
			result = callbackValue.Call([]reflect.Value{result, each, reflect.ValueOf(i)})[0]
		}
	})

	return result.Interface()
}

// Reject function is the opposite of Filter(); This method returns the elements of collection that predicate does not return truthy for.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Reject(data, callback interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, callback)
		if *err != nil {
			return nil
		}

		callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return nil
		}

		validateFuncOutputOneVarBool(err, callbackType, true)
		if *err != nil {
			return nil
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
			if !res[0].Bool() {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Remove function removes all elements from `data` that `predicate` returns truthy for and returns a slice of the removed elements.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data      // type: slice, description: the slice to inspect
//  predicate // type: func(each anyType, i int)bool, description: the function invoked per iteration.
//
// Return values
//
// This function return three values:
//  slice // description: returns slice after elements removed as per `predicate`
//  slice // description: returns slice of removed elements as per `predicate`
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Remove(data interface{}, predicate interface{}) (interface{}, interface{}, error) {
	var result, removed interface{}
	var err error

	func(err *error) {
		defer catch(err)

		result, removed = data, nil

		if !isNonNilData(err, "data", data) {
			return
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)
		removed = makeSlice(dataType).Interface()

		if !isSlice(err, "data", dataValue) {
			return
		}

		callbackValue, callbackType := inspectFunc(err, predicate)
		if *err != nil {
			return
		}

		callbackTypeNumIn := validateFuncInputForSliceLoop(err, callbackType, dataValue)
		if *err != nil {
			return
		}

		validateFuncOutputOneVarBool(err, callbackType, true)
		if *err != nil {
			return
		}

		resultSlice := makeSlice(dataType)
		removedSlice := makeSlice(dataType)

		if dataValueLen == 0 {
			return
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			res := callFuncSliceLoop(callbackValue, each, i, callbackTypeNumIn)
			if res[0].Bool() {
				removedSlice = reflect.Append(removedSlice, each)
			} else {
				resultSlice = reflect.Append(resultSlice, each)
			}
		})

		result = resultSlice.Interface()
		removed = removedSlice.Interface()

		return
	}(&err)

	return result, removed, err
}

// Reverse function reverses `data` so that the first element becomes the last, the second element becomes the second to last, and so on.
//
// Parameters
//
// This function requires one mandatory parameter:
//  data // type: slice, description: the slice to modify
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Reverse(data interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		for i := 0; i < dataValueLen; i++ {
			result = reflect.Append(result, dataValue.Index(dataValueLen-1-i))
		}

		return result.Interface()
	}(&err)

	return result, err
}

// Sample function gets a random element from collection.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Sample(data interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return makeSlice(dataType).Interface()
		}

		return dataValue.Index(RandomInt(0, dataValueLen-1)).Interface()
	}(&err)

	return result, err
}

// SampleSize function gets n random elements at unique keys from collection up to the size of collection.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func SampleSize(data interface{}, take int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if !isPositiveNumber(err, "size", take) {
			return nil
		}

		cache := make(map[int]bool, 0)
		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		if take >= dataValueLen {
			return data
		}

		for result.Len() < take {
			n := RandomInt(0, dataValueLen-1)
			if _, ok := cache[n]; ok {
				continue
			}

			cache[n] = true
			result = reflect.Append(result, dataValue.Index(n))
		}

		return result.Interface()
	}(&err)

	return result, err
}

// Shuffle function creates an array of shuffled values, using a version of the Fisher-Yates shuffle.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Shuffle(data interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, _, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		rand.Seed(time.Now().UnixNano())

		n := dataValueLen
		for i := n - 1; i > 0; i-- {
			j := rand.Intn(i + 1)

			iValue, jValue := dataValue.Index(i).Interface(), dataValue.Index(j).Interface()
			dataValue.Index(i).Set(reflect.ValueOf(jValue))
			dataValue.Index(j).Set(reflect.ValueOf(iValue))
		}

		return dataValue.Interface()
	}(&err)

	return result, err
}

// Size function gets the size of collection by returning its length for array-like values or the number of own enumerable string keyed properties for objects.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to process
//  size // type: number, description: the length of each chunk
//
// Return values
//
// This function return two values:
//  slice // description: returns the new slice of chunks
//  error // description: hold error message if there is an error
//
// Examples
//
// N examples available:
func Size(data interface{}) (int, error) {
	var err error

	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return 0
		}

		dataValue, _, dataValueKind, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.String {
				*err = nil
				return len(data.(string))
			} else if dataValueKind == reflect.Map {
				*err = nil
				return dataValueLen
			}

			return 0
		}

		return dataValueLen
	}(&err)

	return result, err
}

// SortBy is alias of OrderBy()
func SortBy(data, callback interface{}, args ...bool) (interface{}, error) {
	return OrderBy(data, callback, args...)
}

// Tail function gets all but the first element of `data`.
//
// Parameters
//
// This function requires one mandatory parameter:
//  data // type: slice, description: the slice to modify
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Tail(data interface{}) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return makeSlice(dataType).Interface()
		}

		result := dataValue.Slice(1, dataValueLen)
		return result.Interface()
	}(&err)

	return result, err
}

// Take function creates a slice of `data` with `size` elements taken from the beginning.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to query
//  size // type: number, description: the number of elements to take
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Take(data interface{}, size int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if !isZeroOrPositiveNumber(err, "size", size) {
			return data
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			if i < size {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// TakeRight function creates a slice of `data` with `size` elements taken from the end.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to query
//  size // type: number, description: the number of elements to take
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func TakeRight(data interface{}, size int) (interface{}, error) {
	var err error

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if !isZeroOrPositiveNumber(err, "size", size) {
			return data
		}

		result := makeSlice(dataType)

		if dataValueLen == 0 {
			return result.Interface()
		}

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			if i >= (dataValueLen - size) {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)

	return result, err
}

// Union function combines all slices presented on the parameters, then create slice of unique values from it. All slice must have same data type.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data   // type: slice, description: the slice to inspect
//  slice1 // optional, type: slice, description: the slice to inspect
//  slice2 // optional, type: slice, description: the slice to inspect
//  slice3 // optional, type: slice, description: the slice to inspect
//  ...
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Union(data interface{}, slices ...interface{}) (interface{}, error) {
	var err error

	result := (func(err *error) interface{} {
		defer catchWithCustomErrorMessage(err, func(errorMessage string) string {
			if strings.Contains(errorMessage, "is not assignable") {
				return "data type of each elements between slice must be same"
			}

			return errorMessage
		})

		if !isNonNilData(err, "data", data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		result := makeSlice(dataType)
		resultMap := make(map[interface{}]bool, 0)

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			eachRealValue := each.Interface()

			if _, ok := resultMap[eachRealValue]; !ok {
				resultMap[eachRealValue] = true
				result = reflect.Append(result, each)
			}
		})

		for _, each := range slices {
			targetValue, _, _, targetValueLen := inspectData(each)

			if !isSlice(err, "data", targetValue) {
				return nil
			}

			if targetValueLen > 0 {

				forEachSliceStoppable(targetValue, targetValueLen, func(inner reflect.Value, j int) bool {
					if *err != nil {
						return false
					}

					targetEachRealValue := inner.Interface()

					if _, ok := resultMap[targetEachRealValue]; !ok {
						resultMap[targetEachRealValue] = true
						result = reflect.Append(result, inner)
					}

					return true
				})
			}
		}

		return result.Interface()
	})(&err)

	return result, err
}

// Uniq function is same like `Union` but only accept one parameter.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data // type: slice, description: the slice to inspect
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 2 examples available:
func Uniq(data interface{}) (interface{}, error) {
	return Union(data)
}

// Without creates a slice from `data` excluding all given values presented on the `items`.
//
// Parameters
//
// This function requires two mandatory parameters:
//  data  // type: slice, description: the slice to inspect
//  item1 // optional, type: anyType, description: item to exclude
//  item2 // optional, type: anyType, description: item to exclude
//  item3 // optional, type: anyType, description: item to exclude
//  ...
//
// Return values
//
// This function return two values:
//  slice // description: returns slice
//  error // description: hold error message if there is an error
//
// Examples
//
// 3 examples available:
func Without(data interface{}, items ...interface{}) (interface{}, error) {
	return Pull(data, items...)
}
