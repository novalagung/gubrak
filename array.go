package gubrak

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/*
 * Creates an array of elements split into groups the length of size.
 * If array can't be split evenly, the final chunk will be the remaining elements.
 */

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

/*
 * Creates an array with all falsey values removed.
 * The values false, nil, 0, "", (*string)(nil), and other nil values are falsey.
 */

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

/*
 * Creates a new array concatenating array with any additional arrays and/or values.
 */

func Concat(data interface{}, concatenableData ...interface{}) (interface{}, error) {
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

		for i, eachConcatenableData := range concatenableData {
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

/*
 * Creates an array of array values not included in the other given arrays.
 * The order and references of result values are determined by the first array.
 */

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

/*
 * Creates a slice of array with n elements dropped from the beginning.
 */

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

/*
 * Creates a slice of array with n elements dropped from the last.
 */

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

/*
 * Fills elements of array with value from start up to, but not including, end.
 */

func Fill(data, fill interface{}, args ...int) (interface{}, error) {
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

		fillValue, fillType, _, _ := inspectData(fill)

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

/*
 * This method is like _.find except that it returns
 * the index of the first element predicate returns truthy for instead of the element itself.
 */

func FindIndex(data, callback interface{}, args ...int) (int, error) {
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

		callbackValue, callbackType := inspectFunc(err, callback)
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

/*
 * This method is like FindIndex except that it iterates over elements of collection from right to left.
 */

func FindLastIndex(data, callback interface{}, args ...int) (int, error) {
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

		callbackValue, callbackType := inspectFunc(err, callback)
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

/*
 * Gets the first element of array.
 */

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

/*
 * The inverse of _.toPairs; this method returns an object composed from key-value pairs
 */

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

/*
 * Alias of First()
 */

func Head(data interface{}) (interface{}, error) {
	return First(data)
}

/*
 * Gets the index at which the first occurrence of value is found in array using SameValueZero
 * for equality comparisons.
 * If fromIndex is negative, it's used as the offset from the end of array.
 */

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

/*
 * Gets all but the last element of array.
 */

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

/*
 * Creates an array of unique values that are included in all given arrays
 * using SameValueZero for equality comparisons.
 * The order and references of result values are determined by the first array.
 */

func Intersection(data interface{}, compareData ...interface{}) (interface{}, error) {
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

		for _, compare := range compareData {
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

/*
 * Converts all elements in array into a string separated by separator.
 */

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

/*
 * Gets the last element of array.
 */

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

/*
 * This method is like _.indexOf except that it iterates over elements of array from right to left.
 */

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

/*
 * Gets the element at index n of array. If n is negative, the nth element from the end is returned.
 */

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

/*
 * Removes all given values from array using SameValueZero for equality comparisons.
 */

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

/*
 * This method is like _.pull except that it accepts an array of values to remove.
 */

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

/*
 * Removes elements from array corresponding to indexes and returns an array of removed elements.
 */

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

/*
 * Removes all elements from array that predicate returns truthy for
 * and returns an array of the removed elements.
 * The predicate is invoked with three arguments: (value, index, array).
 */

func Remove(data interface{}, callback interface{}) (interface{}, interface{}, error) {
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

		callbackValue, callbackType := inspectFunc(err, callback)
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

/*
 * Reverses array so that the first element becomes the last,
 * the second element becomes the second to last, and so on.
 */

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

/*
 * Gets all but the first element of array.
 */

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

/*
 * Creates a slice of array with n elements taken from the beginning.
 */

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

/*
 * Creates a slice of array with n elements taken from the end.
 */

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

/*
 * Creates an array of unique values, in order,
 * from all given arrays using SameValueZero for equality comparisons.
 */

func Union(data interface{}, target ...interface{}) (interface{}, error) {
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

		for _, each := range target {
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

/*
 * Creates a duplicate-free version of an array,
 * using SameValueZero for equality comparisons,
 * in which only the first occurrence of each element is kept.
 * The order of result values is determined by the order
 * they occur in the array.
 */

func Uniq(data interface{}) (interface{}, error) {
	return Union(data)
}

/*
 * Creates an array excluding all given values using SameValueZero for equality comparisons.
 */

func Without(data interface{}, target ...interface{}) (interface{}, error) {
	return Pull(data, target...)
}

/*
 * This method is like _.xor except that it accepts iteratee which is invoked
 * for each element of each arrays to generate the criterion by which by which they're compared.
 * The order of result values is determined by the order they occur in the arrays.
 * The iteratee is invoked with one argument: (value).
 */

func Xor(data interface{}, compareData ...interface{}) (interface{}, error) {
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

		type ResultMap struct {
			ReflectValue reflect.Value
			counter      int
		}

		resultMap := make(map[interface{}]ResultMap, 0)

		forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
			eachActualValue := each.Interface()
			if val, ok := resultMap[eachActualValue]; !ok {
				resultMap[eachActualValue] = ResultMap{ReflectValue: each, counter: 1}
			} else {
				resultMap[eachActualValue] = ResultMap{ReflectValue: each, counter: val.counter + 1}
			}
		})

		for i, eachCompareData := range compareData {
			eachLabel := fmt.Sprintf("concat data %d", (i + 1))
			eachValue, eachType, _, eachValueLen := inspectData(eachCompareData)

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
				eachActualValue := each.Interface()
				if val, ok := resultMap[eachActualValue]; !ok {
					resultMap[eachActualValue] = ResultMap{ReflectValue: each, counter: 1}
				} else {
					resultMap[eachActualValue] = ResultMap{ReflectValue: each, counter: val.counter + 1}
				}
			})
		}

		result := makeSlice(dataType)

		for _, val := range resultMap {
			if val.counter == 1 {
				result = reflect.Append(result, val.ReflectValue)
			}
		}

		return result.Interface()
	}(&err)

	return result, err
}
