package gubrak

import (
	"fmt"
	"reflect"
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
