package gubrak

import (
	"errors"
	"fmt"
	"reflect"
)

func inspectFunc(err *error, data interface{}) (reflect.Value, reflect.Type) {
	var dataValue reflect.Value
	var dataValueType reflect.Type

	dataValue = reflect.ValueOf(data)

	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	if dataValue.Kind() != reflect.Func {
		*err = errors.New("callback should be function")
		return dataValue, dataValueType
	}

	dataValueType = dataValue.Type()
	return dataValue, dataValueType
}

func inspectData(data interface{}) (reflect.Value, reflect.Type, reflect.Kind, int) {
	var dataValue reflect.Value
	var dataValueType reflect.Type
	var dataValueKind reflect.Kind
	dataValueLen := 0

	if data != nil {
		dataValue = reflect.ValueOf(data)
		dataValueType = dataValue.Type()
		dataValueKind = dataValue.Kind()

		if dataValueKind == reflect.Ptr {
			dataValue = dataValue.Elem()
		}

		if dataValueKind == reflect.Slice || dataValueKind == reflect.Array {
			dataValueLen = dataValue.Len()
		} else if dataValueKind == reflect.Map {
			dataValueLen = len(dataValue.MapKeys())
		}
	}

	return dataValue, dataValueType, dataValueKind, dataValueLen
}

func makeSlice(valueType reflect.Type, args ...int) reflect.Value {
	sliceLen := 0
	sliceCap := 0

	if len(args) > 0 {
		sliceLen = args[0]

		if len(args) > 1 {
			sliceCap = args[1]
		}
	}

	if valueType.Kind() == reflect.Array {
		sliceUnaddresable := reflect.MakeSlice(reflect.SliceOf(valueType.Elem()), 0, 0)
		sliceAddressable := reflect.New(sliceUnaddresable.Type())
		sliceAddressable.Elem().Set(sliceUnaddresable)
		return sliceAddressable.Elem()
	}

	return reflect.MakeSlice(valueType, sliceLen, sliceCap)
}

func validateFuncInputForSliceLoop(err *error, funcType reflect.Type, data reflect.Value) int {
	funcTypeNumIn := funcType.NumIn()

	if funcTypeNumIn == 0 || funcTypeNumIn >= 3 {
		*err = errors.New("callback must only have one or two parameters")
		return funcTypeNumIn
	}

	if funcType.In(0).Kind() != data.Type().Elem().Kind() {
		*err = errors.New("callback 1st parameter's data type should be same with slice element data type")
		return funcTypeNumIn
	}

	if funcTypeNumIn == 2 {
		if funcType.In(1).Kind() != reflect.Int {
			*err = errors.New("callback 2nd parameter's data type should be int")
			return funcTypeNumIn
		}
	}

	return funcTypeNumIn
}

func validateFuncInputForSliceLoopWithoutIndex(err *error, funcType reflect.Type, data reflect.Value) {
	if funcType.NumIn() != 1 {
		*err = errors.New("callback must only have one parameters")
		return
	}

	if funcType.In(0).Kind() != data.Type().Elem().Kind() {
		*err = errors.New("callback parameter's data type should be same with slice data type")
	}
}

func validateFuncInputForCollectionLoop(err *error, funcType reflect.Type, data reflect.Value) int {
	funcTypeNumIn := funcType.NumIn()

	if funcTypeNumIn == 0 || funcTypeNumIn >= 3 {
		*err = errors.New("callback must only have one or two parameters")
		return funcTypeNumIn
	}

	if funcType.In(0).Kind() != data.Type().Elem().Kind() {
		*err = errors.New("callback 1st parameter's data type should be same with map value data type")
		return funcTypeNumIn
	}

	if funcTypeNumIn == 2 {
		if funcType.In(1).Kind() != data.Type().Key().Kind() {
			*err = errors.New("callback 2nd parameter's data type should be same with map key type")
			return funcTypeNumIn
		}
	}

	return funcTypeNumIn
}

func validateFuncOutputNone(err *error, funcType reflect.Type) {
	callbackTypeNumOut := funcType.NumOut()

	if callbackTypeNumOut != 0 {
		*err = errors.New("callback should not have return value")
	}
}

func validateFuncOutputOneVarDynamic(err *error, funcType reflect.Type) int {
	callbackTypeNumOut := funcType.NumOut()
	if callbackTypeNumOut != 1 {
		*err = errors.New("callback return value should only be 1 variable")
		return callbackTypeNumOut
	}

	return callbackTypeNumOut
}

func validateFuncOutputOneVarBool(err *error, callbackType reflect.Type, isMust bool) int {
	isOptional := !isMust

	message := "callback return value should be one variable with bool type"
	if isOptional {
		message = "callback return value data type should be bool, ... or no return value at all"
	}

	callbackTypeNumOut := callbackType.NumOut()
	if callbackTypeNumOut == 1 {
		if callbackType.Out(0).Kind() != reflect.Bool {
			*err = errors.New(message)
			return callbackTypeNumOut
		}
	} else {
		if isOptional {
			if callbackTypeNumOut > 1 {
				*err = errors.New(message)
				return callbackTypeNumOut
			}
		} else {
			*err = errors.New(message)
			return callbackTypeNumOut
		}
	}

	return callbackTypeNumOut
}

func forEachSlice(slice reflect.Value, sliceLen int, eachCallback func(reflect.Value, int)) {
	forEachSliceStoppable(slice, sliceLen, func(each reflect.Value, i int) bool {
		eachDataValue := slice.Index(i)
		eachCallback(eachDataValue, i)
		return true
	})
}

func forEachSliceStoppable(slice reflect.Value, sliceLen int, eachCallback func(reflect.Value, int) bool) {
	for i := 0; i < sliceLen; i++ {
		eachDataValue := slice.Index(i)
		shouldContinue := eachCallback(eachDataValue, i)

		if !shouldContinue {
			return
		}
	}
}

func forEachCollection(collection reflect.Value, keys []reflect.Value, eachCallback func(reflect.Value, reflect.Value, int)) {
	forEachCollectionStoppable(collection, keys, func(value, key reflect.Value, i int) bool {
		eachCallback(value, key, i)
		return true
	})
}

func forEachCollectionStoppable(collection reflect.Value, keys []reflect.Value, eachCallback func(reflect.Value, reflect.Value, int) bool) {
	for i, key := range keys {
		shouldContinue := eachCallback(collection.MapIndex(key), key, i)

		if !shouldContinue {
			return
		}
	}
}

func callFuncSliceLoop(funcToCall, param reflect.Value, i int, numIn int) []reflect.Value {
	if numIn == 1 {
		return funcToCall.Call([]reflect.Value{param})
	}

	return funcToCall.Call([]reflect.Value{param, reflect.ValueOf(i)})
}

func callFuncCollectionLoop(funcToCall, value, key reflect.Value, numIn int) []reflect.Value {
	if numIn == 1 {
		return funcToCall.Call([]reflect.Value{value})
	}

	return funcToCall.Call([]reflect.Value{value, key})
}

func isSlice(err *error, label string, dataValue ...reflect.Value) bool {
	if len(dataValue) == 0 {
		*err = fmt.Errorf("%s cannot be empty", label)
		return false

	} else if len(dataValue) == 1 {
		if dataValue[0].Kind() == reflect.Slice || dataValue[0].Kind() == reflect.Array {
			return true
		}

		*err = fmt.Errorf("%s must be slice", label)
		return false
	} else {
		res := dataValue[0].Kind() == reflect.Slice || dataValue[0].Kind() == reflect.Array

		for i, each := range dataValue {
			if i > 0 {
				res = res || (each.Kind() == reflect.Slice) || (each.Kind() == reflect.Array)
			}
		}

		return res
	}
}

func isNonNilData(err *error, label string, data interface{}) bool {
	if data == nil {
		*err = fmt.Errorf("%s cannot be nil", label)
		return false
	}

	valueOfData := reflect.ValueOf(data)

	switch valueOfData.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer, reflect.Struct:
		if valueOfData.IsNil() {
			*err = fmt.Errorf("%s cannot be nil", label)
			return false
		}
	}

	return true
}

func isZeroOrPositiveNumber(err *error, label string, size int) bool {
	if size < 0 {
		*err = fmt.Errorf("%s must not be negative number", label)
		return false
	} else if size == 0 {
		return true
	}

	return true
}

func isPositiveNumber(err *error, label string, size int) bool {
	if size < 0 {
		*err = fmt.Errorf("%s must be positive number", label)
		return false
	} else if size == 0 {
		*err = fmt.Errorf("%s must be positive number", label)
		return false
	}

	return true
}

func isLeftShouldBeGreaterOrEqualThanRight(err *error, labelLeft string, valueLeft int, labelRight string, valueRight int) bool {
	if valueLeft < valueRight {
		*err = fmt.Errorf("%s should be greater than %s", labelLeft, labelRight)
		return false
	}

	return true
}

func isTypeEqual(err *error, labelLeft string, typeLeft reflect.Type, labelRight string, typeRight reflect.Type) bool {
	if typeLeft != typeRight {
		*err = fmt.Errorf("type of %s should be same with type of %s", labelLeft, labelRight)
		return false
	}

	return true
}

func catch(err *error) {
	if r := recover(); r != nil {
		*err = fmt.Errorf("%v", r)
	}
}

func catchWithCustomErrorMessage(err *error, callback func(string) string) {
	if r := recover(); r != nil {
		*err = errors.New(callback(fmt.Sprintf("%v", r)))
	}
}
