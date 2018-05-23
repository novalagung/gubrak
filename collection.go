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

// Count function creates an object composed of keys generated from the results of running each element of collection thru iteratee. The corresponding value of each key is the number of times the key was returned by iteratee. The iteratee is invoked with one argument: (value).
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

// Each function iterates over elements of array and invokes iteratee for each element. The iteratee is invoked with two arguments: (value, index). Iteratee functions may exit iteration early by explicitly returning false.
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

// Filter function iterates over elements of collection, returning an array of all elements predicate returns truthy for. The predicate is invoked with two arguments: (value, index).
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

// FindLast function is like Find() except that it iterates over elements of collection from right to left.
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

// ForEach is alias of Each()
func ForEach(data, callback interface{}) error {
	return Each(data, callback)
}

// ForEachRight is alias of EachRight()
func ForEachRight(data, callback interface{}) error {
	return EachRight(data, callback)
}

// GroupBy function creates an object composed of keys generated from the results of running each element of collection thru iteratee. The order of grouped values is determined by the order they occur in collection. The corresponding value of each key is an array of elements responsible for generating the key. The iteratee is invoked with one argument: (value).
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

// Map function creates an array of values by running each element in collection thru iteratee. The iteratee is invoked with two arguments: (value, index).
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

// Includes function checks if value is in collection. If collection is a string, it's checked for a substring of value, otherwise SameValueZero is used for equality comparisons. If fromIndex is negative, it's used as the offset from the end of collection.
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

// KeyBy function creates an object composed of keys generated from the results of running each element of collection thru iteratee. The corresponding value of each key is the last element responsible for generating the key. The iteratee is invoked with one argument: (value).
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

// OrderBy sort slices. If orders is unspecified, all values are sorted in ascending order. Otherwise, specify an order of "desc" for descending or "asc" for ascending sort order of corresponding values. The algorithm used is merge sort, as per savigo's post on https://sagivo.com/go-sort-faster-4869bdabc670
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

		isAsync := true
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
							v, _ = strconv.ParseFloat(s, bitSize)
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

// Reduce function reduces collection to a value which is the accumulated result of running each element in collection thru iteratee, where each successive invocation is supplied the return value of the previous. If accumulator is not given, the first element of collection is used as the initial value. The iteratee is invoked with four arguments: (accumulator, value, index|key, collection)
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

// Sample function gets a random element from collection.
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
