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
// This function requires single mandatory parameter:
//  size int // ==> description: the length of each chunk
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Chunk(size int) IChainable {
	g.lastOperation = OperationChunk
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Compact function creates a slice with all falsey values removed from the `data`. These values: `false`, `nil`, `0`, `""`, `(*string)(nil)`, and other nil-able types are considered to be falsey.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Compact() IChainable {
	g.lastOperation = OperationCompact
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

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

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				ok = target.Uint() != 0

			case reflect.Float32, reflect.Float64:
				ok = target.Float() != 0

			case reflect.Complex64, reflect.Complex128:
				ok = target.Complex() != 0

			case reflect.String:
				ok = target.String() != ""

			default: // case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
				ok = !target.IsNil()
			}

			if ok {
				result = reflect.Append(result, each)
			}
		})

		return result.Interface()
	}(&err)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Concat function creates a new slice concatenating `data` with any additional slice.
//
// Parameters
//
// This function requires single mandatory parameter:
//  sliceToConcat interface{} // ==> description: the slice to concatenate
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Concat(sliceToConcat interface{}) IChainable {
	g.lastOperation = OperationConcat
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _concat(&err, g.data, sliceToConcat)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// ConcatMany function creates a new slice concatenating `data` with any additional slices (the 2nd parameter and rest).
//
// Parameters
//
// This function requires optional variadic parameters:
//  sliceToConcat1 interface{} // ==> description: the slice to concatenate
//  sliceToConcat2 interface{} // ==> description: the slice to concatenate
//  sliceToConcat3 interface{} // ==> description: the slice to concatenate
//  ...
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) ConcatMany(slicesToConcat ...interface{}) IChainable {
	g.lastOperation = OperationConcatMany
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _concat(&err, g.data, slicesToConcat...)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

func _concat(err *error, data interface{}, slicesToConcat ...interface{}) interface{} {
	defer catch(err)

	if !isNonNilData(err, "data", data) {
		return nil
	}

	dataValue, dataType, _, dataValueLen := inspectData(data)

	if !isSlice(err, "data", dataValue) {
		return nil
	}

	result := makeSlice(dataType)

	forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
		result = reflect.Append(result, each)
	})

	for i, eachConcatenableData := range slicesToConcat {
		eachLabel := fmt.Sprintf("concat data %d", (i + 1))
		eachValue, eachType, _, eachValueLen := inspectData(eachConcatenableData)

		if !isSlice(err, eachLabel, eachValue) {
			return nil
		}

		if !isTypeEqual(err, "data", dataType.Elem(), eachLabel, eachType.Elem()) {
			return nil
		}

		if eachValueLen == 0 {
			continue
		}

		forEachSlice(eachValue, eachValueLen, func(each reflect.Value, i int) {
			result = reflect.Append(result, each)
		})
	}

	return result.Interface()
}

// Contains function checks if value is in data. If data is a string, it's checked for a substring of value, otherwise SameValueZero is used for equality comparisons. If `fromIndex` is negative, it's used as the offset from the end of data.
//
// Parameters
//
// This function requires single mandatory parameter:
//  search interface{} // ==> description: the value to search for.
//  fromIndex int      // ==> optional
//                     //     description: The index to search from
//                     //     default value: 0
//
// Return values
//
// Chain with these methods to get result:
//  .Result() bool                  // ==> description: returns true if value is found, else false
//  .ResultAndError() (bool, error) // ==> description: returns true if value is found, else false, and error object
//  .Error() error                  // ==> description: returns error object
//  .IsError() bool                 // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Contains(search interface{}, args ...int) IChainableBoolResult {
	g.lastOperation = OperationContains
	if g.IsError() || g.shouldReturn() {
		return &resultContains{chainable: g}
	}

	err := (error)(nil)
	result := func(err *error) bool {
		defer catch(err)

		if dataValue, dataOK := g.data.(string); dataOK {
			if searchValue, searchOK := search.(string); searchOK {
				return strings.Contains(dataValue, searchValue)
			}

			return false
		}

		if !isNonNilData(err, "data", g.data) {
			return false
		}

		dataValue, _, dataValueKind, dataValueLen := inspectData(g.data)

		startIndex := 0
		if len(args) > 0 {
			startIndex = args[0]
		}

		if !isZeroOrPositiveNumber(err, "start index", startIndex) {
			return false
		}

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.Map {
				*err = nil
				return _containsCollection(err, dataValue, search, startIndex)
			}

			*err = errors.New((*err).Error() + ", map, or a string")
			return false
		}

		return _containsSlice(err, dataValue, dataValueLen, search, startIndex)
	}(&err)
	if err != nil {
		return &resultContains{chainable: g.markError(result, err)}
	}

	return &resultContains{chainable: g.markResult(result)}
}

func _containsSlice(err *error, dataValue reflect.Value, dataValueLen int, search interface{}, startIndex int) bool {
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

func _containsCollection(err *error, dataValue reflect.Value, search interface{}, startIndex int) bool {
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

// Count get the length of `data`.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() int                  // ==> description: returns length of data
//  .ResultAndError() (int, error) // ==> description: returns length of data, and error object
//  .Error() error                 // ==> description: returns error object
//  .IsError() bool                // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Count() IChainableNumberResult {
	g.lastOperation = OperationCount
	if g.IsError() || g.shouldReturn() {
		return &resultCount{chainable: g}
	}

	err := (error)(nil)
	result := _count(&err, g.data, nil)
	if err != nil {
		return &resultCount{chainable: g.markError(result, err)}
	}

	return &resultCount{chainable: g.markResult(result)}
}

// CountBy get the length of `data` filtered by `iteratee`.
//
// Parameters
//
// This function requires single mandatory parameter:
//  iteratee interface{} // ==> type: `func(each anyType, i int)bool` or
//                       //           `func(value anyType, key anyType, i int)bool`
//                       //     description: the function invoked per iteration
//
// Return values
//
// Chain with these methods to get result:
//  .Result() int                  // ==> description: returns the result after operation
//  .ResultAndError() (int, error) // ==> description: returns the result after operation, and error object
//  .Error() error                 // ==> description: returns error object
//  .IsError() bool                // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) CountBy(iteratee interface{}) IChainableNumberResult {
	g.lastOperation = OperationCountBy
	if g.IsError() || g.shouldReturn() {
		return &resultCount{chainable: g}
	}

	err := (error)(nil)
	result := _count(&err, g.data, iteratee)
	if err != nil {
		return &resultCount{chainable: g.markError(result, err)}
	}

	return &resultCount{chainable: g.markResult(result)}
}

func _count(err *error, data, predicate interface{}) int {
	defer catch(err)

	if !isNonNilData(err, "data", data) {
		return 0
	}

	dataValue, dataValueType, dataValueKind, dataValueLen := inspectData(data)

	if !isSlice(err, "data", dataValue) {
		if dataValueKind == reflect.Map {
			*err = nil
			return _countCollection(err, dataValue, dataValueType, dataValueKind, dataValueLen, predicate)
		}

		return 0
	}

	return _countSlice(err, dataValue, dataValueType, dataValueKind, dataValueLen, predicate)
}

func _countSlice(err *error, dataValue reflect.Value, dataValueType reflect.Type, dataValueKind reflect.Kind, dataValueLen int, callback interface{}) int {

	var callbackValue reflect.Value
	var callbackType reflect.Type
	var callbackTypeNumIn int

	if callback == nil {
		return dataValueLen
	}

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

	if callback == nil {
		return dataValueLen
	}

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

// Difference function creates a slice of `data` that values not included in the other given slice. The order and references of result values are determined by the first slice.
//
// Parameters
//
// This function requires single mandatory parameter:
//  dataToCompare interface{} // ==> description: the slice to differentiate
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Difference(dataToCompare interface{}) IChainable {
	g.lastOperation = OperationDifferenceMany
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _difference(&err, g.data, dataToCompare)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// DifferenceMany function creates a slice of `data` that values not included in the other given slices. The order and references of result values are determined by the first slice.
//
// Parameters
//
// This function requires optional variadic parameters:
//  datasToCompare1 interface{} // ==> description: the slice to differentiate
//  datasToCompare2 interface{} // ==> description: the slice to differentiate
//  datasToCompare3 interface{} // ==> description: the slice to differentiate
//  ...
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) DifferenceMany(datasToCompare ...interface{}) IChainable {
	g.lastOperation = OperationDifferenceMany
	if g.IsError() || g.shouldReturn() {
		return g
	}

	if len(datasToCompare) == 0 {
		return g.markError(nil, errors.New("data to compare cannot be empty"))
	}

	err := (error)(nil)
	result := _difference(&err, g.data, datasToCompare...)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

func _difference(err *error, data interface{}, dataToCompare ...interface{}) interface{} {
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

	dataToCompareMap := make(map[reflect.Value]int)
	for i, each := range dataToCompare {
		eachLabel := fmt.Sprintf("difference data %d", (i + 1))
		eachValue, eachType, _, eachValueLen := inspectData(each)

		if !isSlice(err, eachLabel, eachValue) {
			return nil
		}

		if !isTypeEqual(err, "data", dataType, eachLabel, eachType) {
			return nil
		}

		dataToCompareMap[eachValue] = eachValueLen
	}

	forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
		isFound := false

		for compareValue, compareValueLen := range dataToCompareMap {
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
}

// Drop function creates a slice of `data` with `n` elements dropped from the beginning.
//
// Parameters
//
// This function requires single mandatory parameter:
//  size int // ==> description: the number of elements to drop
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Drop(size int) IChainable {
	g.lastOperation = OperationDrop
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if !isZeroOrPositiveNumber(err, "size", size) {
			return g.data
		}

		if size == 0 {
			return g.data
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// DropRight function creates a slice of `data` with `n` elements dropped from the end.
//
// Parameters
//
// This function requires single mandatory parameter:
//  size int // ==> description: the number of elements to drop
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) DropRight(size int) IChainable {
	g.lastOperation = OperationDropRight
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if !isZeroOrPositiveNumber(err, "size", size) {
			return g.data
		}

		if size == 0 {
			return g.data
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Each iterates over elements of `data` and invokes `iteratee` for each element. Iteratee functions may exit iteration early by explicitly returning false
//
// Parameters
//
// This function requires single mandatory parameter:
//  iteratee interface{} // ==> type: `func(each anyType, i int)` or
//                       //           `func(each anyType, i int)bool` or
//                       //           `func(value anyType, key anyType, i int)` or
//                       //           `func(value anyType, key anyType, i int)bool`
//                       // ==> description: the function invoked per iteration.
//                       //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                       //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                       //                  and both are optional.
//                       //                  if return value is provided then the next iteration is controlled by returned value.
//                       //                  `return true` will make the iteration continue, meanwhile `return false` will stop it
//
// Return values
//
// Chain with these methods to get result:
//  .Error() error  // ==> description: returns error object
//  .IsError() bool // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Each(iteratee interface{}) IChainableNoReturnValueResult {
	g.lastOperation = OperationEach
	if g.IsError() || g.shouldReturn() {
		return &resultEach{chainable: g}
	}

	err := (error)(nil)
	_each(&err, g.data, iteratee, true)
	if err != nil {
		return &resultEach{chainable: g.markError(nil, err)}
	}

	return &resultEach{chainable: g.markResult(nil)}
}

// EachRight iterates over elements of `data` from tail to head, and invokes `iteratee` for each element. Iteratee functions may exit iteration early by explicitly returning false
//
// Parameters
//
// This function requires single mandatory parameter:
//  iteratee interface{} // ==> type: `func(each anyType, i int)` or
//                       //           `func(each anyType, i int)bool` or
//                       //           `func(value anyType, key anyType, i int)` or
//                       //           `func(value anyType, key anyType, i int)bool`
//                       // ==> description: the function invoked per iteration.
//                       //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                       //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                       //                  and both are optional.
//                       //                  if return value is provided then the next iteration is controlled by returned value.
//                       //                  `return true` will make the iteration continue, meanwhile `return false` will stop it
//
// Return values
//
// Chain with these methods to get result:
//  .Error() error  // ==> description: returns error object
//  .IsError() bool // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) EachRight(iteratee interface{}) IChainableNoReturnValueResult {
	g.lastOperation = OperationEachRight
	if g.IsError() || g.shouldReturn() {
		return &resultEach{chainable: g}
	}

	err := (error)(nil)
	_each(&err, g.data, iteratee, true)
	if err != nil {
		return &resultEach{chainable: g.markError(nil, err)}
	}

	return &resultEach{chainable: g.markResult(nil)}
}

func _each(err *error, data, iteratee interface{}, isForward bool) {
	defer catch(err)

	if !isNonNilData(err, "data", data) {
		return
	}

	dataValue, dataValueType, dataValueKind, dataValueLen := inspectData(data)

	if !isSlice(err, "data", dataValue) {
		if dataValueKind == reflect.Map {
			*err = nil
			_eachCollection(err, dataValue, dataValueType, dataValueKind, dataValueLen, iteratee, isForward)
		}

		return
	}

	_eachSlice(err, dataValue, dataValueType, dataValueKind, dataValueLen, iteratee, isForward)
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

// Exclude function removes value from `data`.
//
// Parameters
//
// This function requires single mandatory parameter:
//  itemToExclude interface{} // ==> description: the item to exclude
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Exclude(itemToExclude interface{}) IChainable {
	g.lastOperation = OperationExclude
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _exclude(&err, g.data, itemToExclude)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// ExcludeMany function removes all given values from `data`.
//
// Parameters
//
// This function requires optional variadic parameters:
//  itemToExclude1 interface{} // ==> description: the item to exclude
//  itemToExclude2 interface{} // ==> description: the item to exclude
//  itemToExclude3 interface{} // ==> description: the item to exclude
//  ...
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) ExcludeMany(itemsToExclude ...interface{}) IChainable {
	g.lastOperation = OperationExcludeMany
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _exclude(&err, g.data, itemsToExclude...)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

func _exclude(err *error, data interface{}, items ...interface{}) interface{} {
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
}

// ExcludeAt function removes value by index from `data`.
//
// Parameters
//
// This function requires single mandatory parameter:
//  indexOfItemToExclude interface{} // ==> description: the index of item to exclude
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) ExcludeAt(indexOfItemToExclude int) IChainable {
	g.lastOperation = OperationExcludeAt
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _excludeAt(&err, g.data, indexOfItemToExclude)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// ExcludeAtMany function removes value by index from `data`.
//
// Parameters
//
// This function requires optional variadic parameters:
//  indexOfItemToExclude1 interface{} // ==> description: the index of item to exclude
//  indexOfItemToExclude2 interface{} // ==> description: the index of item to exclude
//  indexOfItemToExclude3 interface{} // ==> description: the index of item to exclude
//  ...
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) ExcludeAtMany(indexesOfItemToExclude ...int) IChainable {
	g.lastOperation = OperationExcludeAtMany
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _excludeAt(&err, g.data, indexesOfItemToExclude...)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

func _excludeAt(err *error, data interface{}, indexes ...int) interface{} {
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
}

// Fill function fills elements of `data` with `value` from `start` up to, but not including, `end`.
//
// Parameters
//
// This function requires single mandatory parameter, `value`; and two other optional parameters:
//  value interface{} // ==> description: the value to fill slice with. This variable's data type must be same with slice's element data type
//  start int         // ==> optional
//                    //     description: the start position
//                    //     default value: 0
//  end int           // ==> optional
//                    //     description: the end position
//                    //     default value: len(data)
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Fill(value interface{}, args ...int) IChainable {
	g.lastOperation = OperationFill
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

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
			return g.data
		}

		if !isZeroOrPositiveNumber(err, "last index", lastIndex) {
			return g.data
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Filter function iterates over elements of slice or struct object or map, returning an array of all elements predicate returns truthy for.
//
// Parameters
//
// This function requires single mandatory parameter:
//  predicate interface{} // ==> type: `func(each anyType, i int)bool` or
//                        //           `func(value anyType, key anyType, i int)bool`
//                        // ==> description: the function invoked per iteration.
//                        //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                        //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                        //                  and both are optional.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Filter(predicate interface{}) IChainable {
	g.lastOperation = OperationFilter
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, dataValueKind, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.Map {
				*err = nil
				return _filterCollection(err, dataValue, dataType, dataValueKind, dataValueLen, predicate)
			}

			return nil
		}

		return _filterSlice(err, dataValue, dataType, dataValueKind, dataValueLen, predicate)
	}(&err)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
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

// Find function iterates over elements of collection, returning the first element predicate returns truthy for.
//
// Parameters
//
// This function requires single mandatory parameter:
//  predicate interface{} // ==> type: `func(each anyType, i int)bool` or
//                        //           `func(value anyType, key anyType, i int)bool`
//                        // ==> description: the function invoked per iteration.
//                        //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                        //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                        //                  and both are optional.
//  fromIndex int         // ==> optional
//                        //     description: The index to search from
//                        //     default value: 0
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Find(predicate interface{}, args ...int) IChainable {
	g.lastOperation = OperationFind
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, predicate)
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// FindIndex function iterates over elements of collection, returning the index of first element predicate returns truthy for.
//
// Parameters
//
// This function requires single mandatory parameter:
//  predicate interface{} // ==> type: `func(each anyType, i int)bool` or
//                        //           `func(value anyType, key anyType, i int)bool`
//                        // ==> description: the function invoked per iteration.
//                        //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                        //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                        //                  and both are optional.
//  fromIndex int         // ==> optional
//                        //     description: The index to search from
//                        //     default value: 0
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) FindIndex(predicate interface{}, args ...int) IChainable {
	g.lastOperation = OperationFindIndex
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return -1
		}

		dataValue, _, _, dataValueLen := inspectData(g.data)

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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// FindLast function iterates over elements from tail to head, returning the first element predicate returns truthy for.
//
// Parameters
//
// This function requires single mandatory parameter:
//  predicate interface{} // ==> type: `func(each anyType, i int)bool` or
//                        //           `func(value anyType, key anyType, i int)bool`
//                        // ==> description: the function invoked per iteration.
//                        //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                        //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                        //                  and both are optional.
//  fromIndex int         // ==> optional
//                        //     description: The index to search from
//                        //     default value: len(data)-1
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) FindLast(predicate interface{}, args ...int) IChainable {
	g.lastOperation = OperationFindLast
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, predicate)
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// FindLastIndex function iterates over elements from tail to head, returning the index of first element predicate returns truthy for.
//
// Parameters
//
// This function requires single mandatory parameter:
//  predicate interface{} // ==> type: `func(each anyType, i int)bool` or
//                        //           `func(value anyType, key anyType, i int)bool`
//                        // ==> description: the function invoked per iteration.
//                        //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                        //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                        //                  and both are optional.
//  fromIndex int         // ==> optional
//                        //     description: The index to search from
//                        //     default value: len(data)-1
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) FindLastIndex(predicate interface{}, args ...int) IChainable {
	g.lastOperation = OperationFindLastIndex
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return -1
		}

		dataValue, _, _, dataValueLen := inspectData(g.data)

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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// First function gets the first element of `data`.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) First() IChainable {
	g.lastOperation = OperationFirst
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, _, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return nil
		}

		return dataValue.Index(0).Interface()
	}(&err)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// FromPairs function returns an object composed from key-value `data`.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) FromPairs() IChainable {
	g.lastOperation = OperationFromPairs
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataValueType, _, dataValueLen := inspectData(g.data)

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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// GroupBy function creates an object composed of keys generated from the results of running each element of collection thru iteratee. The order of grouped values is determined by the order they occur in collection. The corresponding value of each key is an array of elements responsible for generating the key.
//
// Parameters
//
// This function requires single mandatory parameter:
//  predicate interface{} // ==> type: `func(each anyType, i int)<any type>`
//                        // ==> description: the function invoked per iteration.
//                        //                  the 2nd argument represents index of each element, and it's optional.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) GroupBy(predicate interface{}) IChainable {
	g.lastOperation = OperationGroupBy
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, predicate)
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// IndexOf function gets the index at which the first occurrence of `search` is found in `data`. If `fromIndex` is negative, it's used as the offset from the end of `data`.
//
// Parameters
//
// This function requires single mandatory parameter:
//  search interface{} // ==> description: the value to search for.
//  fromIndex int      // ==> optional
//                     //     description: The index to search from
//                     //     default value: 0
//
// Return values
//
// Chain with these methods to get result:
//  .Result() int                  // ==> description: return the index of found element
//  .ResultAndError() (int, error) // ==> description: return the index of found element, and error object
//  .Error() error                 // ==> description: returns error object
//  .IsError() bool                // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) IndexOf(search interface{}, args ...int) IChainableNumberResult {
	g.lastOperation = OperationIndexOf
	if g.IsError() || g.shouldReturn() {
		return &resultIndexOf{chainable: g}
	}

	err := (error)(nil)
	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return -1
		}

		dataValue, _, _, dataValueLen := inspectData(g.data)

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
	if err != nil {
		return &resultIndexOf{chainable: g.markError(result, err)}
	}

	return &resultIndexOf{chainable: g.markResult(result)}
}

// Initial function gets all but the last element of `data`.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Initial() IChainable {
	g.lastOperation = OperationInitial
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return makeSlice(dataType).Interface()
		}

		return dataValue.Slice(0, dataValueLen-1).Interface()
	}(&err)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Intersection function creates a slice of unique values that are included in all given slice. The order and references of result values are determined by the first slice.
//
// Parameters
//
// This function requires single mandatory parameter:
//  dataToIntersect interface{} // ==> description: the slice to intersect
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Intersection(dataIntersect interface{}) IChainable {
	g.lastOperation = OperationIntersection
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _intersection(&err, g.data, dataIntersect)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// IntersectionMany function creates a slice of unique values that are included in all given slices. The order and references of result values are determined by the first slice.
//
// Parameters
//
// This function requires optional variadic parameters:
//  dataToIntersect1 interface{} // ==> description: the slice to intersect
//  dataToIntersect2 interface{} // ==> description: the slice to intersect
//  dataToIntersect3 interface{} // ==> description: the slice to intersect
//  ...
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) IntersectionMany(dataToIntersects ...interface{}) IChainable {
	g.lastOperation = OperationIntersection
	if g.IsError() || g.shouldReturn() {
		return g
	}

	if len(dataToIntersects) == 0 {
		return g.markError(nil, errors.New("data intersects cannot be nil"))
	}

	err := (error)(nil)
	result := _intersection(&err, g.data, dataToIntersects...)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

func _intersection(err *error, data interface{}, dataIntersects ...interface{}) interface{} {
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
}

// Join function converts all elements in `data` into a string separated by `separator`.
//
// Parameters
//
// This function requires single mandatory parameter:
//  separator string // ==> description: the element joiner
//
// Return values
//
// Chain with these methods to get result:
//  .Result() string                  // ==> description: returns the result after operation
//  .ResultAndError() (string, error) // ==> description: returns the result after operation, and error object
//  .Error() error                    // ==> description: returns error object
//  .IsError() bool                   // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Join(separator string) IChainableStringResult {
	g.lastOperation = OperationJoin
	if g.IsError() || g.shouldReturn() {
		return &resultJoin{chainable: g}
	}

	err := (error)(nil)
	result := func(err *error) string {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return ""
		}

		dataValue, _, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return ""
		}

		if val, ok := g.data.([]string); ok {
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
	if err != nil {
		return &resultJoin{chainable: g.markError(result, err)}
	}

	return &resultJoin{chainable: g.markResult(result)}
}

// KeyBy function creates an object composed of keys generated from the results of running each element of collection thru iteratee. The corresponding value of each key is the last element responsible for generating the key.
//
// Parameters
//
// This function requires single mandatory parameter:
//  predicate interface{} // ==> type: `func(each anyType, i int)<any type>` or
//                        //           `func(value anyType, key anyType, i int)<any type>`
//                        // ==> description: the function invoked per iteration.
//                        //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                        //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                        //                  and both are optional.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) KeyBy(predicate interface{}) IChainable {
	g.lastOperation = OperationKeyBy
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataValueType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, predicate)
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Last function gets the last element of `data`.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Last() IChainable {
	g.lastOperation = OperationLast
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, _, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return nil
		}

		return dataValue.Index(dataValueLen - 1).Interface()
	}(&err)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// LastIndexOf function iterates the element from tail to head, then return the index at which the first occurrence of `search` is found in `data`. If `fromIndex` is negative, it's used as the offset from the end of `data`.
//
// Parameters
//
// This function requires single mandatory parameter:
//  search interface{} // ==> description: the value to search for.
//  fromIndex int      // ==> optional
//                     //     description: The index to search from
//                     //     default value: len(data)-1
//
// Return values
//
// Chain with these methods to get result:
//  .Result() int                  // ==> description: return the index of found element
//  .ResultAndError() (int, error) // ==> description: return the index of found element, and error object
//  .Error() error                 // ==> description: returns error object
//  .IsError() bool                // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) LastIndexOf(search interface{}, args ...int) IChainableNumberResult {
	g.lastOperation = OperationLast
	if g.IsError() || g.shouldReturn() {
		return &resultLastIndexOf{chainable: g}
	}

	err := (error)(nil)
	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return -1
		}

		dataValue, _, _, dataValueLen := inspectData(g.data)

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
	if err != nil {
		return &resultLastIndexOf{chainable: g.markError(result, err)}
	}

	return &resultLastIndexOf{chainable: g.markResult(result)}
}

// Map function creates an array of values by running each element in `data` thru iteratee.
//
// Parameters
//
// This function requires single mandatory parameter:
//  callback interface{} // ==> type: `func(each anyType, i int)<any type>` or
//                       //           `func(value anyType, key anyType, i int)<any type>`
//                       // ==> description: the function invoked per iteration.
//                       //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                       //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                       //                  and both are optional.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Map(callback interface{}) IChainable {
	g.lastOperation = OperationMap
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, _, _, dataValueLen := inspectData(g.data)

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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Nth function gets the element at index `n` of `data`. If `n` is negative, the nth element from the end is returned.
//
// Parameters
//
// This function requires single mandatory parameter:
//  index int // ==> description: the index of the element to return
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Nth(index int) IChainable {
	g.lastOperation = OperationNth
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)

	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, _, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return nil
		}

		if index < 0 {
			index = dataValueLen + index
		}

		if index < dataValueLen {
			return dataValue.Index(index).Interface()
		}

		return nil
	}(&err)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// OrderBy sort slices. If orders is unspecified, all values are sorted in ascending order. Otherwise, specify an order of "desc" for descending or "asc" for ascending sort order of corresponding values. The algorithm used is merge sort, as per savigo's post on https://sagivo.com/go-sort-faster-4869bdabc670
//
// Parameters
//
// This function requires single mandatory parameter:
//  predicate interface{} // ==> type: `func(each anyType, i int)<any type>` or
//                        //           `func(value anyType, key anyType, i int)<any type>`
//                        // ==> description: the function invoked per iteration.
//                        //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                        //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                        //                  and both are optional.
//  isAscending bool      // ==> optional
//                        //     description: the sort order. `true` for ascending, and `false` for descending.
//                        //     default value: true
//  isAsync bool          // ==> optional
//                        //     description: concurrent sort. set to `true` to enable pararel sorting (faster for certain data structure)
//                        //     default value: false
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) OrderBy(predicate interface{}, args ...bool) IChainable {
	g.lastOperation = OperationOrderBy
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _orderBy(&err, g.data, predicate, args...)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

func _orderBy(err *error, data, callback interface{}, args ...bool) interface{} {
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
}

// Partition function creates an array of elements split into two groups, the first of which contains elements predicate returns truthy for, the second of which contains elements predicate returns falsey for. The predicate is invoked with one argument: (value).
//
// Parameters
//
// This function requires single mandatory parameter:
//  predicate interface{} // ==> type: `func(each anyType, i int)bool`
//                        // ==> description: the function invoked per iteration.
//                        //                  the 2nd argument represents index of each element, and it's optional.
//                        //                  and both are optional.
//
// Return values
//
// Chain with these methods to get result:
//  .ResultTruthy() interface{}                         // ==> description: return slice of elements which predicate returns truthy for
//  .ResultFalsey() interface{}                         // ==> description: return slice of elements which predicate returns falsey for
//  .ResultAndError() (interface{}, interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                                      // ==> description: returns error object
//  .IsError() bool                                     // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Partition(callback interface{}) IChainableTwoReturnValueResult {
	g.lastOperation = OperationPartition
	if g.IsError() || g.shouldReturn() {
		return &resultPartition{chainable: g}
	}

	err := (error)(nil)
	truhty, falsey := func(err *error) (interface{}, interface{}) {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil, nil
		}

		dataValue, dataValueType, _, dataValueLen := inspectData(g.data)

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
	if err != nil {
		return &resultPartition{chainable: g.markError([]interface{}{truhty, falsey}, err)}
	}

	return &resultPartition{chainable: g.markResult([]interface{}{truhty, falsey})}
}

// Reduce function reduces collection to a value which is the accumulated result of running each element in collection thru iteratee, where each successive invocation is supplied the return value of the previous. If accumulator is not given, the first element of collection is used as the initial value.
//
// Parameters
//
// This function require two mandatory parameters:
//  iteratee interface{} // ==> type: `func(accumulator <any type>, each anyType, i int)<any type>` or
//                       //           `func(accumulator <any type>, value anyType, key anyType, i int)<any type>`
//                       // ==> description: the function invoked per iteration.
//                       //                  the 1st argument is the accumulator. at first the value is coming from `initial`
//                       //                  for slice, the 3rd argument represents index of each element, and it's optional.
//                       //                  for struct object/map, the 3rd and 4th arguments represent key and index of each item respectively,
//                       //                  and both are optional.
//  initial interface{}  // ==> description: the initial value.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Reduce(iteratee, initial interface{}) IChainable {
	g.lastOperation = OperationReduce
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataValueType, dataValueKind, dataValueLen := inspectData(g.data)
		if dataValueLen == 0 {
			return initial
		}

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.Map {
				*err = nil
				return _reduceCollection(err, dataValue, dataValueType, dataValueKind, dataValueLen, iteratee, initial)
			}

			return nil
		}

		return _reduceSlice(err, dataValue, dataValueType, dataValueKind, dataValueLen, iteratee, initial)
	}(&err)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
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

// Reject function iterates over elements of slice or struct object or map, returning an array of all elements predicate returns FALSEY for.
//
// Parameters
//
// This function requires single mandatory parameter:
//  predicate interface{} // ==> type: `func(each anyType, i int)bool` or
//                        //           `func(value anyType, key anyType, i int)bool`
//                        // ==> description: the function invoked per iteration.
//                        //                  for slice, the 2nd argument represents index of each element, and it's optional.
//                        //                  for struct object/map, the 2nd and 3rd arguments represent key and index of each item respectively,
//                        //                  and both are optional.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Reject(predicate interface{}) IChainable {
	g.lastOperation = OperationReject
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		callbackValue, callbackType := inspectFunc(err, predicate)
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Reverse function reverses `data` so that the first element becomes the last, the second element becomes the second to last, and so on.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Reverse() IChainable {
	g.lastOperation = OperationReverse
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Sample function gets a random element from `data`.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Sample() IChainable {
	g.lastOperation = OperationSample
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return makeSlice(dataType).Interface()
		}

		return dataValue.Index(RandomInt(0, dataValueLen-1)).Interface()
	}(&err)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// SampleSize function gets slice of random elements from `data`.
//
// Parameters
//
// This function requires single mandatory parameter:
//  take int // ==> description: the length of each chunk
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) SampleSize(take int) IChainable {
	g.lastOperation = OperationSampleSize
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

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
			return g.data
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Shuffle function creates a slice of shuffled values, using a version of the Fisher-Yates shuffle.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Shuffle() IChainable {
	g.lastOperation = OperationShuffle
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, _, _, dataValueLen := inspectData(g.data)

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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Size function gets the size of slice or struct object/map by returning its length for array-like values or the number of own enumerable string keyed properties for objects.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() int                  // ==> description: returns the length of data
//  .ResultAndError() (int, error) // ==> description: returns the length of data, and error object
//  .Error() error                 // ==> description: returns error object
//  .IsError() bool                // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Size() IChainable {
	g.lastOperation = OperationSize
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) int {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return 0
		}

		dataValue, _, dataValueKind, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			if dataValueKind == reflect.String {
				*err = nil
				return len(g.data.(string))
			} else if dataValueKind == reflect.Map {
				*err = nil
				return dataValueLen
			}

			return 0
		}

		return dataValueLen
	}(&err)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Tail function gets all but the first element of `data`.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() int                  // ==> description: returns the length of data
//  .ResultAndError() (int, error) // ==> description: returns the length of data, and error object
//  .Error() error                 // ==> description: returns error object
//  .IsError() bool                // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Tail() IChainable {
	g.lastOperation = OperationTail
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if dataValueLen == 0 {
			return makeSlice(dataType).Interface()
		}

		result := dataValue.Slice(1, dataValueLen)
		return result.Interface()
	}(&err)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Take function creates a slice of `data` with `size` elements taken from the beginning.
//
// Parameters
//
// This function requires single mandatory parameter:
//  size int // ==> description: the length of each chunk
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Take(size int) IChainable {
	g.lastOperation = OperationTake
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if !isZeroOrPositiveNumber(err, "size", size) {
			return g.data
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// TakeRight function creates a slice of `data` with `size` elements taken from the end.
//
// Parameters
//
// This function requires single mandatory parameter:
//  size int // ==> description: the length of each chunk
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) TakeRight(size int) IChainable {
	g.lastOperation = OperationTakeRight
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := func(err *error) interface{} {
		defer catch(err)

		if !isNonNilData(err, "data", g.data) {
			return nil
		}

		dataValue, dataType, _, dataValueLen := inspectData(g.data)

		if !isSlice(err, "data", dataValue) {
			return nil
		}

		if !isZeroOrPositiveNumber(err, "size", size) {
			return g.data
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Uniq create slice of unique values from it.
//
// Parameters
//
// This function does not requires any parameter.
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) Uniq() IChainable {
	g.lastOperation = OperationUniq
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _union(&err, g.data)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// Union function combines all slices presented on the parameters, then create slice of unique values from it. All slice must have same data type.
//
// Parameters
//
// This function requires optional variadic parameters:
//  sliceToUnion1 interface{} // ==> description: the index of item to exclude
//  sliceToUnion2 interface{} // ==> description: the index of item to exclude
//  sliceToUnion3 interface{} // ==> description: the index of item to exclude
//  ...
//
// Return values
//
// Chain with these methods to get result:
//  .Result() interface{}                  // ==> description: returns the result after operation
//  .ResultAndError() (interface{}, error) // ==> description: returns the result after operation, and error object
//  .Error() error                         // ==> description: returns error object
//  .IsError() bool                        // ==> description: return `true` on error, otherwise `false`
//
// Examples
//
// List of examples available:
func (g *Chainable) UnionMany(sliceToUnion ...interface{}) IChainable {
	g.lastOperation = OperationUnionMany
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _union(&err, g.data, sliceToUnion...)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

func _union(err *error, data interface{}, slices ...interface{}) interface{} {
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
}
