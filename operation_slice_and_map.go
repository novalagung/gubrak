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

// ============================================== Chunk

const OperationChunk = "Chunk()"

type IChainableChunk interface {
	Chunk(int) IChainable
}

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

// ============================================== Compact

const OperationCompact = "Compact()"

type IChainableCompact interface {
	Compact() IChainable
}

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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// ============================================== Concat

const OperationConcatMany = "ConcatMany()"
const OperationConcat = "Concat()"

type IChainableConcat interface {
	ConcatMany(...interface{}) IChainable
	Concat(interface{}) IChainable
}

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

	if dataValueLen == 0 {
		return result.Interface()
	}

	forEachSlice(dataValue, dataValueLen, func(each reflect.Value, i int) {
		result = reflect.Append(result, each)
	})

	for i, eachConcatenableData := range slicesToConcat {
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
}

// ============================================== Count

const (
	OperationCountBy = "CountBy()"
	OperationCount   = "Count()"
)

type IChainableCount interface {
	CountBy(interface{}) IChainableCountResult
	Count() IChainableCountResult
}

type IChainableCountResult interface {
	ResultAndError() (int, error)
	Result() int
	Error() error
	IsError() bool
}

type resultCount struct {
	chainable *Chainable
	IChainableCountResult
}

func (g *resultCount) ResultAndError() (int, error) {
	return g.Result(), g.Error()
}

func (g *resultCount) Result() int {
	v, _ := g.chainable.data.(int)
	return v
}

func (g *resultCount) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultCount) IsError() bool {
	return g.Error() != nil
}

func (g *Chainable) CountBy(predicate interface{}) IChainableCountResult {
	g.lastOperation = OperationCountBy
	if g.IsError() || g.shouldReturn() {
		return &resultCount{chainable: g}
	}

	err := (error)(nil)
	result := _count(&err, g.data, predicate)
	if err != nil {
		return &resultCount{chainable: g.markError(result, err)}
	}

	return &resultCount{chainable: g.markResult(result)}
}

func (g *Chainable) Count() IChainableCountResult {
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

// ============================================== Difference

const (
	OperationDifferenceMany = "DifferenceMany()"
	OperationDifference     = "Difference()"
)

type IChainableDifference interface {
	DifferenceMany(...interface{}) IChainable
	Difference(interface{}) IChainable
}

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
			continue
		}

		if dataValueLen == 0 {
			continue
		}

		if !isTypeEqual(err, "data", dataType, eachLabel, eachType) {
			continue
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

// ============================================== Drop

const (
	OperationDrop      = "Drop()"
	OperationDropRight = "DropRight()"
)

type IChainableDrop interface {
	Drop(int) IChainable
	DropRight(int) IChainable
}

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

// ============================================== Each

const (
	OperationEach         = "Each()"
	OperationEachRight    = "EachRight()"
	OperationForEach      = "ForEach()"
	OperationForEachRight = "ForEachRight()"
)

type IChainableEach interface {
	Each(interface{}) IChainableEachResult
	EachRight(interface{}) IChainableEachResult
	ForEach(interface{}) IChainableEachResult
	ForEachRight(interface{}) IChainableEachResult
}

type IChainableEachResult interface {
	Error() error
	IsError() bool
}

type resultEach struct {
	chainable *Chainable
	IChainableCountResult
}

func (g *resultEach) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultEach) IsError() bool {
	return g.Error() != nil
}

func (g *Chainable) Each(callback interface{}) IChainableEachResult {
	g.lastOperation = OperationEach
	if g.IsError() || g.shouldReturn() {
		return &resultEach{chainable: g}
	}

	err := (error)(nil)
	_each(&err, g.data, callback, true)
	if err != nil {
		return &resultEach{chainable: g.markError(nil, err)}
	}

	return &resultEach{chainable: g.markResult(nil)}
}

func (g *Chainable) EachRight(callback interface{}) IChainableEachResult {
	g.lastOperation = OperationEachRight
	if g.IsError() || g.shouldReturn() {
		return &resultEach{chainable: g}
	}

	err := (error)(nil)
	_each(&err, g.data, callback, true)
	if err != nil {
		return &resultEach{chainable: g.markError(nil, err)}
	}

	return &resultEach{chainable: g.markResult(nil)}
}

func (g *Chainable) ForEach(callback interface{}) IChainableEachResult {
	return g.Each(callback)
}

func (g *Chainable) ForEachRight(callback interface{}) IChainableEachResult {
	return g.EachRight(callback)
}

func _each(err *error, data, callback interface{}, isForward bool) {
	defer catch(err)

	if !isNonNilData(err, "data", data) {
		return
	}

	dataValue, dataValueType, dataValueKind, dataValueLen := inspectData(data)

	if !isSlice(err, "data", dataValue) {
		if dataValueKind == reflect.Map {
			*err = nil
			_eachCollection(err, dataValue, dataValueType, dataValueKind, dataValueLen, callback, isForward)
		}

		return
	}

	_eachSlice(err, dataValue, dataValueType, dataValueKind, dataValueLen, callback, isForward)
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

// ============================================== Fill

const (
	OperationFill = "Fill()"
)

type IChainableFill interface {
	Fill(interface{}, ...int) IChainable
}

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

// ============================================== Filter

const (
	OperationFilter = "Filter()"
)

type IChainableFilter interface {
	Filter(interface{}) IChainable
}

func (g *Chainable) Filter(callback interface{}) IChainable {
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
				return _filterCollection(err, dataValue, dataType, dataValueKind, dataValueLen, callback)
			}

			return nil
		}

		return _filterSlice(err, dataValue, dataType, dataValueKind, dataValueLen, callback)
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

// ============================================== Find

const (
	OperationFind          = "Find()"
	OperationFindIndex     = "FindIndex()"
	OperationFindLast      = "FindLast()"
	OperationFindLastIndex = "FindLastIndex()"
)

type IChainableFind interface {
	Find(interface{}, ...int) IChainable
	FindIndex(interface{}, ...int) IChainable
	FindLast(interface{}, ...int) IChainable
	FindLastIndex(interface{}, ...int) IChainable
}

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

// ============================================== First

const (
	OperationFirst = "First()"
	OperationHead  = "Head()"
)

type IChainableFirst interface {
	First() IChainable
	Head() IChainable
}

func (g *Chainable) First() IChainable {
	g.lastOperation = OperationFirst
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _first(&err, g.data)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

func (g *Chainable) Head() IChainable {
	g.lastOperation = OperationHead
	if g.IsError() || g.shouldReturn() {
		return g
	}

	err := (error)(nil)
	result := _first(&err, g.data)
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

func _first(err *error, data interface{}) interface{} {
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
}

// ============================================== FromPairs

const (
	OperationFromPairs = "FromPairs()"
)

type IChainableFromPairs interface {
	FromPairs() IChainable
}

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

// ============================================== GroupBy

const (
	OperationGroupBy = "GroupBy()"
)

type IChainableGroupBy interface {
	GroupBy(interface{}) IChainable
}

func (g *Chainable) GroupBy(callback interface{}) IChainable {
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// ============================================== Includes

const (
	OperationIncludes = "Includes()"
)

type IChainableIncludes interface {
	Includes(interface{}, ...int) IChainableIncludesResult
}

type IChainableIncludesResult interface {
	ResultAndError() (bool, error)
	Result() bool
	Error() error
	IsError() bool
}

type resultIncludes struct {
	chainable *Chainable
	IChainableIncludesResult
}

func (g *resultIncludes) ResultAndError() (bool, error) {
	return g.Result(), g.Error()
}

func (g *resultIncludes) Result() bool {
	v, _ := g.chainable.data.(bool)
	return v
}

func (g *resultIncludes) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultIncludes) IsError() bool {
	return g.Error() != nil
}

func (g *Chainable) Includes(search interface{}, args ...int) IChainableIncludesResult {
	g.lastOperation = OperationIncludes
	if g.IsError() || g.shouldReturn() {
		return &resultIncludes{chainable: g}
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
	if err != nil {
		return &resultIncludes{chainable: g.markError(result, err)}
	}

	return &resultIncludes{chainable: g.markResult(result)}
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

// ============================================== IndexOf

const (
	OperationIndexOf = "IndexOf()"
)

type IChainableIndexOf interface {
	IndexOf(interface{}, ...int) IChainableIndexOfResult
}

type IChainableIndexOfResult interface {
	ResultAndError() (int, error)
	Result() int
	Error() error
	IsError() bool
}

type resultIndexOf struct {
	chainable *Chainable
	IChainableIncludesResult
}

func (g *resultIndexOf) ResultAndError() (int, error) {
	return g.Result(), g.Error()
}

func (g *resultIndexOf) Result() int {
	v, _ := g.chainable.data.(int)
	return v
}

func (g *resultIndexOf) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultIndexOf) IsError() bool {
	return g.Error() != nil
}

func (g *Chainable) IndexOf(search interface{}, args ...int) IChainableIndexOfResult {
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

// ============================================== Initial

const (
	OperationInitial = "Initial()"
)

type IChainableInitial interface {
	Initial(interface{}) IChainable
}

func (g *Chainable) Initial(callback interface{}) IChainable {
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

// ============================================== Intersection

const (
	OperationIntersection     = "Intersection()"
	OperationIntersectionMany = "IntersectionMany()"
)

type IChainableIntersection interface {
	Intersection(interface{}) IChainable
	IntersectionMany(data ...interface{}) IChainable
}

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

func (g *Chainable) IntersectionMany(dataIntersects ...interface{}) IChainable {
	g.lastOperation = OperationIntersection
	if g.IsError() || g.shouldReturn() {
		return g
	}

	if len(dataIntersects) == 0 {
		return g.markError(nil, errors.New("data intersects cannot be nil"))
	}

	err := (error)(nil)
	result := _intersection(&err, g.data, dataIntersects...)
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

// ============================================== Join

const (
	OperationJoin = "Join()"
)

type IChainableJoin interface {
	Join(string) IChainableJoinResult
}

type IChainableJoinResult interface {
	ResultAndError() (string, error)
	Result() string
	Error() error
	IsError() bool
}

type resultJoin struct {
	chainable *Chainable
	IChainableIncludesResult
}

func (g *resultJoin) ResultAndError() (string, error) {
	return g.Result(), g.Error()
}

func (g *resultJoin) Result() string {
	v, _ := g.chainable.data.(string)
	return v
}

func (g *resultJoin) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultJoin) IsError() bool {
	return g.Error() != nil
}

func (g *Chainable) Join(separator string) IChainableJoinResult {
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

// ============================================== KeyBy

const (
	OperationKeyBy = "KeyBy()"
)

type IChainableKeyBy interface {
	KeyBy(interface{}) IChainable
}

func (g *Chainable) KeyBy(callback interface{}) IChainable {
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
	if err != nil {
		return g.markError(result, err)
	}

	return g.markResult(result)
}

// ================================================================ RAW

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

func SortBy(data, callback interface{}, args ...bool) (interface{}, error) {
	return OrderBy(data, callback, args...)
}

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

func Uniq(data interface{}) (interface{}, error) {
	return Union(data)
}

func Without(data interface{}, items ...interface{}) (interface{}, error) {
	return Pull(data, items...)
}
