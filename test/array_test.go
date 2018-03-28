package test

import (
	. "github.com/novalagung/gubrak"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestChunkNegativeSize(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	result, err := Chunk(data, -1)
	// ===> []

	assert.NotNil(t, err)
	assert.EqualError(t, err, "size must not be negative number")
	assert.Nil(t, result)
}

func TestChunkZeroSize(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	result, err := Chunk(data, 0)
	// ===> []

	assert.Nil(t, err)
	assert.Equal(t, make([][]string, 0), result)
}

func TestChunkSizeTwoInt(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	result, err := Chunk(data, 2)
	// ===> [][]int{ {1, 2}, {3, 4}, {5} }

	assert.Nil(t, err)
	assert.EqualValues(t, [][]int{{1, 2}, {3, 4}, {5}}, result)
}

func TestChunkSizeThreeString(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e"}
	result, err := Chunk(data, 3)
	// ===> [][]string{ {"a", "b", "c"}, {"d", "e"} }

	assert.Nil(t, err)
	assert.EqualValues(t, [][]string{{"a", "b", "c"}, {"d", "e"}}, result)
}

func TestChunkSliceInterface(t *testing.T) {
	data := []interface{}{
		3.2,
		"a",
		-1,
		make([]byte, 0),
		map[string]int{"b": 2},
		[]string{"a", "b", "c"},
	}
	result, err := Chunk(data, 2)
	// ===> [][]interface{}{
	//        {3.2, "a"},
	//        {-1, []uint8{}},
	//        {map[string]int{"b": 2}, []string{"a", "b", "c"}},
	//      }

	assert.Nil(t, err)
	assert.EqualValues(t, [][]interface{}{
		{3.2, "a"},
		{-1, []uint8{}},
		{map[string]int{"b": 2}, []string{"a", "b", "c"}},
	}, result)
}

func TestChunkSizeTwo(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	result, err := Chunk(data, 2)
	// ===> [["a", "b"], ["c", "d"]]

	assert.Nil(t, err)
	assert.EqualValues(t, [][]string{{"a", "b"}, {"c", "d"}}, result)
}

func TestChunkSizeThree(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	result, err := Chunk(data, 3)
	// ===> [["a", "b", "c"], ["d"]]

	assert.Nil(t, err)
	assert.EqualValues(t, [][]string{{"a", "b", "c"}, {"d"}}, result)
}

func TestChunkSizeFour(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	result, err := Chunk(data, 4)
	// ===> [["a", "b", "c", "d"]]

	assert.Nil(t, err)
	assert.EqualValues(t, [][]string{{"a", "b", "c", "d"}}, result)
}

func TestChunkSizeAHundred(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	result, err := Chunk(data, 1000)
	// ===> [["a", "b", "c", "d"]]

	assert.Nil(t, err)
	assert.EqualValues(t, [][]string{{"a", "b", "c", "d"}}, result)
}

func TestChunkNilData(t *testing.T) {
	result, err := Chunk(nil, 2)
	// ===> nil

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data cannot be nil")
	assert.Nil(t, result)
}

func TestChunkEmptyData(t *testing.T) {
	data := []string{}
	result, err := Chunk(data, 2)
	// ===> []

	assert.Nil(t, err)
	assert.EqualValues(t, [][]string{}, result)
}

func TestChunkStringData(t *testing.T) {
	data := "hello"
	result, err := Chunk(data, 2)
	// ===> nil

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data must be slice")
	assert.Nil(t, result)
}

func TestChunkIntData(t *testing.T) {
	data := 12
	result, err := Chunk(data, 2)
	// ===> nil

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data must be slice")
	assert.Nil(t, result)
}

func TestChunkPointerData(t *testing.T) {
	data := "hello"
	result, err := Chunk(&data, 2)
	// ===> nil

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data must be slice")
	assert.Nil(t, result)
}

func TestChunkMapData(t *testing.T) {
	data := map[string]string{
		"fruit":     "manggo",
		"vegetable": "spinach",
	}
	result, err := Chunk(&data, 2)
	// ===> nil

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data must be slice")
	assert.Nil(t, result)
}

func TestCompactFewData(t *testing.T) {
	var dataInterfaceNil interface{} = nil
	var dataInterface interface{} = "damian"
	var dataPointerStringNil *string = nil
	var dataPointerString *string = (func(s string) *string { return &s })("damian")
	dataSliceIntEmpty := make([]int, 0)
	dataSliceInt := []int{1, 2, 3}
	dataMapEmpty := make(map[string]string, 0)
	dataMap := map[string]string{"name": "damian"}

	data := []interface{}{
		-2, 0, 1, 2,
		false, true,
		"", "damian",
		nil,
		dataInterfaceNil,
		dataInterface,
		dataPointerStringNil,
		dataPointerString,
		dataSliceIntEmpty,
		dataSliceInt,
		dataMapEmpty,
		dataMap,
	}

	result, err := Compact(data)
	resultParsed := result.([]interface{})
	/* ===> [
		 0: -2,
		 1: 1,
		 2: 2,
		 3: true,
		 4: "damian",
		 5: "damian",
		 6: (*string)(0xc42008f2b0),
		 7: []int{},
		 8: []int{1, 2, 3},
		 9: map[string]string{},
		10: map[string]string{"name":"damian"}
	] */

	assert.Empty(t, err)
	assert.Len(t, resultParsed, 11)

	assert.Equal(t, -2, resultParsed[0])
	assert.Equal(t, 1, resultParsed[1])
	assert.Equal(t, 2, resultParsed[2])
	assert.Equal(t, true, resultParsed[3])
	assert.Equal(t, "damian", resultParsed[4])
	assert.Equal(t, "damian", resultParsed[5])
	assert.Equal(t, "damian", *(resultParsed[6].(*string)))

	assert.Equal(t, make([]int, 0), resultParsed[7])
	assert.Len(t, resultParsed[7], 0)

	assert.EqualValues(t, []int{1, 2, 3}, resultParsed[8])
	assert.Len(t, resultParsed[8], 3)

	assert.Equal(t, make(map[string]string, 0), resultParsed[9])
	assert.Len(t, resultParsed[9], 0)

	assert.Equal(t, map[string]string{"name": "damian"}, resultParsed[10])
	assert.Equal(t, "damian", resultParsed[10].(map[string]string)["name"])
	assert.Len(t, resultParsed[10], 1)
}

func TestCompactInt(t *testing.T) {
	data := []int{-2, -1, 0, 1, 2}

	result, err := Compact(data)
	resultParsed := result.([]int)
	/* ===> [-2, -1, 1, 2] */

	assert.Empty(t, err)
	assert.Len(t, resultParsed, 4)

	assert.Equal(t, []int{-2, -1, 1, 2}, resultParsed)
}

func TestCompactString(t *testing.T) {
	data := []string{"a", "b", "", "d"}

	result, err := Compact(data)
	resultParsed := result.([]string)
	/* ===> []string{"a", "b", "d"} */

	assert.Empty(t, err)
	assert.Len(t, resultParsed, 3)

	assert.Equal(t, []string{"a", "b", "d"}, resultParsed)
}

func TestCompactPointerString(t *testing.T) {
	item1, item2, item3 := "a", "b", "c"
	data := []*string{&item1, nil, &item2, nil, &item3}

	result, err := Compact(data)
	resultParsed := result.([]*string)
	/* ===> []string{"a", "b", "d"} */

	assert.Empty(t, err)
	assert.Len(t, resultParsed, 3)
}

func TestCompactNilData(t *testing.T) {
	result, err := Compact(nil)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data cannot be nil")
	assert.Nil(t, result)
}

func TestConcatIntData(t *testing.T) {
	data := []int{1, 2, 3, 4}
	data1 := []int{4, 6, 7}
	data2 := []int{8, 9}
	result, err := Concat(data, data1, data2)
	// ===> []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 3, 4, 4, 6, 7, 8, 9}, result)
}

func TestConcatStringData(t *testing.T) {
	data := []string{"my"}
	dataConcat1 := []string{"name", "is"}
	dataConcat2 := []string{"jason", "todd"}
	result, err := Concat(data, dataConcat1, dataConcat2)
	// ===> []string{ "my", "name", "is", "jason", "todd" }

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"my", "name", "is", "jason", "todd"}, result)
}

func TestConcatNilData(t *testing.T) {
	result, err := Concat(nil)
	// ===> nil

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data cannot be nil")
	assert.Nil(t, result)
}

func TestConcatWithNil(t *testing.T) {
	data := []int{1, 2, 3, 4}
	result, err := Concat(data, nil)
	// ===> []int{1, 2, 3, 4}

	assert.NotNil(t, err)
	assert.EqualError(t, err, "concat data 1 must be slice")
	assert.EqualValues(t, []int{1, 2, 3, 4}, result)
}

func TestDifferenceOneData(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 6, 7}
	diff := []int{2, 7}
	result, err := Difference(data, diff)
	// ===> []int{1, 3, 4, 4, 6}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 3, 4, 4, 6}, result)
}

func TestDifferenceOneMultipleData(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 6, 7}
	diff1 := []int{2, 7}
	diff2 := []int{1, 2, 3}
	diff3 := []int{4, 7}
	result, err := Difference(data, diff1, diff2, diff3)
	// ===> []int{6}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{6}, result)
}

func TestDifferenceStringData(t *testing.T) {
	data := []string{"a", "b", "b", "c", "d", "e", "f", "g", "h"}
	dataDiff1 := []string{"b", "d"}
	dataDiff2 := []string{"e", "f", "h"}
	result, err := Difference(data, dataDiff1, dataDiff2)
	// ===> []string{ "a", "c", "g" }

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"a", "c", "g"}, result)
}

func TestDifferenceFloatData(t *testing.T) {
	data := []float64{1.1, 1.11, 1.2, 2.3, 3.0, 3, 4.0, 4.00000, 4.000000001}
	dataDiff1 := []float64{1.1, 3}
	dataDiff2 := []float64{4.000000001}
	result, err := Difference(data, dataDiff1, dataDiff2)
	// ===> []float64{ 1.11, 1.2, 2.3, 4, 4 }

	assert.Nil(t, err)
	assert.EqualValues(t, []float64{1.11, 1.2, 2.3, 4, 4}, result)
}

func TestDifferenceNilData(t *testing.T) {
	diff := []int{2, 7}
	result, err := Difference(nil, diff)
	// ===> nil

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data cannot be nil")
	assert.Nil(t, result)
}

func TestDifferenceWithNilDiffData(t *testing.T) {
	data := []int{2, 7}
	result, err := Difference(data, nil)
	// ===> []int{2, 7}

	assert.NotNil(t, err)
	assert.EqualError(t, err, "difference data 1 must be slice")
	assert.EqualValues(t, []int{2, 7}, result)
}

func TestDropZeroSize(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	size := 0
	result, err := Drop(data, size)
	// ===> []int{1, 2, 3, 4, 4, 5, 6}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 3, 4, 4, 5, 6}, result)
}

func TestDropZeroSizeOne(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	size := 1
	result, err := Drop(data, size)
	// ===> []int{2, 3, 4, 4, 5, 6}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{2, 3, 4, 4, 5, 6}, result)
}

func TestDropZeroSizeThree(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	size := 3
	result, err := Drop(data, size)
	// ===> []int{4, 4, 5, 6}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{4, 4, 5, 6}, result)
}

func TestDropZeroSizeTen(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	size := 10
	result, err := Drop(data, size)
	// ===> []int{}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{}, result)
}

func TestDropZeroSizeNegative(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	size := -2
	result, err := Drop(data, size)
	// ===> []int{}

	assert.NotNil(t, err)
	assert.EqualError(t, err, "size must not be negative number")
	assert.EqualValues(t, []int{1, 2, 3, 4, 4, 5, 6}, result)
}

func TestDropRightZeroSize(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	size := 0
	result, err := DropRight(data, size)
	// ===> []int{1, 2, 3, 4, 4, 5, 6}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 3, 4, 4, 5, 6}, result)
}

func TestDropRightZeroSizeOne(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	size := 1
	result, err := DropRight(data, size)
	// ===> []int{1, 2, 3, 4, 4, 5}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 3, 4, 4, 5}, result)
}

func TestDropRightZeroSizeThree(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	size := 3
	result, err := DropRight(data, size)
	// ===> []int{1, 2, 3, 4}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 3, 4}, result)
}

func TestDropRightZeroSizeTen(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	size := 10
	result, err := DropRight(data, size)
	// ===> []int{}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{}, result)
}

func TestDropRightZeroSizeNegative(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	size := -2
	result, err := DropRight(data, size)
	// ===> []int{}

	assert.NotNil(t, err)
	assert.EqualError(t, err, "size must not be negative number")
	assert.EqualValues(t, []int{1, 2, 3, 4, 4, 5, 6}, result)
}

func TestFillWithANumber(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	replacement := 9
	result, err := Fill(data, replacement)
	// ===> []int{9, 9, 9, 9, 9, 9, 9}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{9, 9, 9, 9, 9, 9, 9}, result)
}

func TestFillWithStartIndex(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	replacement := 9
	startIndex := 2
	result, err := Fill(data, replacement, startIndex)
	// ===> []int{1, 2, 9, 9, 9, 9, 9}

	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 9, 9, 9, 9, 9}, result)
}

func TestFillWithAStringAndStartIndex(t *testing.T) {
	data := []string{"grayson", "jason", "tim", "damian"}
	replacement := "alfred"
	startIndex := 2
	result, err := Fill(data, replacement, startIndex)
	// ===> []int{"grayson", "jason", "alfred", "alfred"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"grayson", "jason", "alfred", "alfred"}, result)
}

func TestFillWithStartIndexAndLastIndex(t *testing.T) {
	data := []float64{1, 2.2, 3.0002, 4, 4, 5.12, 6}
	replacement := float64(9)
	startIndex := 3
	lastIndex := 5
	result, err := Fill(data, replacement, startIndex, lastIndex)
	// ===> []float64{1, 2.2, 3.0002, 9, 9, 5.12, 6}

	assert.Nil(t, err)
	assert.EqualValues(t, []float64{1, 2.2, 3.0002, 9, 9, 5.12, 6}, result)
}

func TestFillWithNegativeStartIndex(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	replacement := 9
	startIndex := -1
	result, err := Fill(data, replacement, startIndex)
	// ===> []int{1, 2, 3, 4, 4, 5, 6}

	assert.NotNil(t, err)
	assert.EqualError(t, err, "start index must not be negative number")
	assert.EqualValues(t, []int{1, 2, 3, 4, 4, 5, 6}, result)
}

func TestFillWithNegativeLastIndex(t *testing.T) {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	replacement := 9
	startIndex := 2
	lastIndex := -1
	result, err := Fill(data, replacement, startIndex, lastIndex)
	// ===> []int{1, 2, 3, 4, 4, 5, 6}

	assert.NotNil(t, err)
	assert.EqualError(t, err, "last index must not be negative number")
	assert.EqualValues(t, []int{1, 2, 3, 4, 4, 5, 6}, result)
}

func TestFindIndex(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := FindIndex(data, func(each string) bool {
		return each == "tim"
	})
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestFindIndexFloat64(t *testing.T) {
	data := []float64{1, 1.1, 1.2, 1.200001, 1.2000000001, 1.3}
	result, err := FindIndex(data, func(each float64) bool {
		return each == 1.2000000001
	})
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestFindIndexWithWrongData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := FindIndex(data, func(each string) bool {
		return each == "hello"
	})
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestFindIndexWithWrongCallbackReturnType(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := FindIndex(data, func(each string) int {
		return 12
	})
	// ===> -1

	assert.NotNil(t, err)
	assert.EqualError(t, err, "callback return value should be one variable with bool type")
	assert.Equal(t, -1, result)
}

func TestFindIndexWithWrongCallbackParameter(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := FindIndex(data, func(each int) bool {
		return each == 12
	})
	// ===> -1

	assert.NotNil(t, err)
	assert.EqualError(t, err, "callback 1st parameter's data type should be same with slice element data type")
	assert.Equal(t, -1, result)
}

func TestFindIndexWithFromIndex(t *testing.T) {
	data := []int{1, 2, 3, 3, 4, 5}
	result, err := FindIndex(data, func(each int) bool {
		return each == 3
	}, 2)
	// ===> 2

	assert.Nil(t, err)
	assert.Equal(t, 2, result)
}

func TestFindIndexWithAnotherFromIndex(t *testing.T) {
	data := []int{1, 2, 3, 3, 4, 5}
	result, err := FindIndex(data, func(each int) bool {
		return each == 3
	}, 3)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestFindLastIndex(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := FindLastIndex(data, func(each string) bool {
		return each == "tim"
	})
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestFindLastIndexWithWrongData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := FindLastIndex(data, func(each string) bool {
		return each == "hello"
	})
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestFindLastIndexWithWrongCallbackReturnType(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := FindLastIndex(data, func(each string) int {
		return 12
	})
	// ===> -1

	assert.NotNil(t, err)
	assert.EqualError(t, err, "callback return value should be one variable with bool type")
	assert.Equal(t, -1, result)
}

func TestFindLastIndexWithWrongCallbackParameter(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := FindLastIndex(data, func(each int) bool {
		return each == 12
	})
	// ===> -1

	assert.NotNil(t, err)
	assert.EqualError(t, err, "callback 1st parameter's data type should be same with slice element data type")
	assert.Equal(t, -1, result)
}

func TestLastFindIndexWithFromIndex(t *testing.T) {
	data := []int{1, 2, 2, 3, 3, 4, 5}
	result, err := FindLastIndex(data, func(each int) bool {
		return each == 3
	})
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestLastFindIndexWithAnotherFromIndex(t *testing.T) {
	data := []int{1, 2, 2, 3, 3, 4, 5}
	result, err := FindLastIndex(data, func(each int) bool {
		return each == 3
	}, 3)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestLastFindIndexWithYetAnotherFromIndex(t *testing.T) {
	data := []int{1, 2, 2, 3, 3, 4, 5}
	result, err := FindLastIndex(data, func(each int) bool {
		return each == 3
	}, 2)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestFirst(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := First(data)
	// ===> "damian"

	assert.Nil(t, err)
	assert.Equal(t, "damian", result)
}

func TestFirstWithEmptyData(t *testing.T) {
	data := []string{}
	result, err := First(data)
	// ===> nil

	assert.Nil(t, err)
	assert.Equal(t, nil, result)
}

func TestFirstWithNilData(t *testing.T) {
	result, err := First(nil)
	// ===> nil

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data cannot be nil")
	assert.Equal(t, nil, result)
}

func TestFromPairsWithDataStringInt(t *testing.T) {
	data := []interface{}{
		[]interface{}{"a", 1},
		[]interface{}{"b", 2},
	}
	result, err := FromPairs(data)

	assert.Nil(t, err)
	assert.Equal(t, map[interface{}]interface{}{
		"a": 1,
		"b": 2,
	}, result)
}

func TestFromPairsWithDataBoolSlice(t *testing.T) {
	data := []interface{}{
		[]interface{}{true, []int{1, 2, 3}},
		[]interface{}{false, []string{"damian", "grayson"}},
	}
	result, err := FromPairs(data)

	assert.Nil(t, err)
	assert.Equal(t, map[interface{}]interface{}{
		true:  []int{1, 2, 3},
		false: []string{"damian", "grayson"},
	}, result)
}

func TestFromPairsWithDataSliceInterface(t *testing.T) {
	data := []interface{}{
		[]interface{}{[]int{1, 2, 3}, "hello"},
		[]interface{}{[]string{"damian", "grayson"}, "hello"},
	}
	result, err := FromPairs(data)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "runtime error: hash of unhashable type []int")
	assert.Equal(t, nil, result)
}

func TestFromPairsWithInvalidType(t *testing.T) {
	data := []string{"a", "b", "c"}
	result, err := FromPairs(data)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "supported type only []interface{}")
	assert.Nil(t, result)
}

func TestIndexOf(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := IndexOf(data, "tim")
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestIndexOfWithFromIndex(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := IndexOf(data, "tim", 4)
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestIndexOfWithWrongSearchData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := IndexOf(data, "hello")
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestIndexOfWithWrongSearchDataWithDifferentType(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := IndexOf(data, make(map[string]string, 0))
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestIndexOfWithFromIndexMinus7(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", -7)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestIndexOfWithFromIndexMinus6(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", -6)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestIndexOfWithFromIndexMinus5(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", -5)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestIndexOfWithFromIndexMinus4(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", -4)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestIndexOfWithFromIndexMinus3(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", -3)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestIndexOfWithFromIndexMinus2(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", -2)
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestIndexOfWithFromIndexMinus1(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", -1)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestIndexOfWithFromIndexZero(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", 0)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestIndexOfWithFromIndex1(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", 1)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestIndexOfWithFromIndex2(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", 2)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestIndexOfWithFromIndex3(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", 3)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestIndexOfWithFromIndex4(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", 4)
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestIndexOfWithFromIndex5(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", 5)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestIndexOfWithFromIndex6(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", 6)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestIndexOfWithFromIndex7(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := IndexOf(data, "tim", 7)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestIndexOfWithInvalidSearchData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := IndexOf(data, 12)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestIndexOfWithNilSearchData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := IndexOf(data, nil)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestInitial(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Initial(data)
	// ===> []string{"damian", "grayson"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "grayson"}, result)
}

func TestInitialEmptyData(t *testing.T) {
	data := []int{}
	result, err := Initial(data)
	// ===> nil

	assert.Nil(t, err)
	assert.EqualValues(t, []int{}, result)
}

func TestIntersection(t *testing.T) {
	result, err := Intersection(
		[]string{"damian", "grayson", "cassandra", "tim", "tim", "jason"},
		[]string{"cassandra", "tim", "jason"},
		[]string{"cassandra", "jason"},
	)
	// ===> []string{"cassandra", "jason"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"cassandra", "jason"}, result)
}

func TestIntersectionWithEmptyComparison(t *testing.T) {
	result, err := Intersection(
		[]string{"damian", "grayson", "cassandra"},
		[]string{},
	)
	// ===> []string{"cassandra", "jason"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{}, result)
}

func TestJoin(t *testing.T) {
	result, err := Join([]string{"damian", "grayson", "cassandra"}, "|")
	// ===> "damian|grayson|cassandra"

	assert.Nil(t, err)
	assert.Equal(t, "damian|grayson|cassandra", result)
}

func TestJoinIntData(t *testing.T) {
	result, err := Join([]int{1, 2, 3, 4}, ",")
	// ===> "1,2,3,4"

	assert.Nil(t, err)
	assert.Equal(t, "1,2,3,4", result)
}

func TestLast(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Last(data)
	// ===> "cassandra"

	assert.Nil(t, err)
	assert.Equal(t, "cassandra", result)
}

func TestLastWithEmptyData(t *testing.T) {
	data := []string{}
	result, err := Last(data)
	// ===> nil

	assert.Nil(t, err)
	assert.Equal(t, nil, result)
}

func TestLastWithNilData(t *testing.T) {
	result, err := Last(nil)
	// ===> nil

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data cannot be nil")
	assert.Equal(t, nil, result)
}

func TestLastIndexOfWithFromIndexMinus7(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", -7)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestLastIndexOfWithFromIndexMinus6(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", -6)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestLastIndexOfWithFromIndexMinus5(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", -5)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestLastIndexOfWithFromIndexMinus4(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", -4)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestLastIndexOfWithFromIndexMinus3(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", -3)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestLastIndexOfWithFromIndexMinus2(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", -2)
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestLastIndexOfWithFromIndexMinus1(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", -1)
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestLastIndexOfWithFromIndexZero(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", 0)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestLastIndexOfWithFromIndex1(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", 1)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestLastIndexOfWithFromIndex2(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", 2)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestLastIndexOfWithFromIndex3(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", 3)
	// ===> 3

	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestLastIndexOfWithFromIndex4(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", 4)
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestLastIndexOfWithFromIndex5(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", 5)
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestLastIndexOfWithFromIndex6(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", 6)
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestLastIndexOfWithFromIndex7(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := LastIndexOf(data, "tim", 7)
	// ===> 4

	assert.Nil(t, err)
	assert.Equal(t, 4, result)
}

func TestLastIndexOfWithWrongSearchData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := LastIndexOf(data, "hello")
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestLastIndexOfWithInvalidSearchData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := LastIndexOf(data, 12)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestLastIndexOfWithNilSearchData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := LastIndexOf(data, nil)
	// ===> -1

	assert.Nil(t, err)
	assert.Equal(t, -1, result)
}

func TestNth(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Nth(data, 1)
	// ===> "grayson"

	assert.Nil(t, err)
	assert.Equal(t, "grayson", result)
}

func TestNthWrongIndex(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Nth(data, 7)
	// ===> "grayson"

	assert.Nil(t, err)
	assert.Equal(t, nil, result)
}

func TestNthNegativeIndex(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Nth(data, -1)
	// ===> "cassandra"

	assert.Nil(t, err)
	assert.Equal(t, "cassandra", result)
}

func TestPullOneData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := Pull(data, "tim")
	// ===> []string{"damian", "grayson", "cassandra", "jason"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "grayson", "cassandra", "jason"}, result)
}

func TestPullThreeData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := Pull(data, "tim", "grayson", "cassandra")
	// ===> []string{"damian", "jason"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "jason"}, result)
}

func TestPullWithNoDataToPull(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := Pull(data)
	// ===> []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}, result)
}

func TestPullAllOneData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := PullAll(data, []string{"tim"})
	// ===> []string{"damian", "grayson", "cassandra", "jason"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "grayson", "cassandra", "jason"}, result)
}

func TestPullAllThreeData(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := PullAll(data, []string{"tim", "grayson", "cassandra"})
	// ===> []string{"damian", "jason"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "jason"}, result)
}

func TestPullAt(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := PullAt(data, 1, 2, 3)
	// ===> []string{"damian", "tim", "jason"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "tim", "jason"}, result)
}

func TestPullAtInvalidIndex(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := PullAt(data, -2, 3)
	// ===> []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}

	assert.NotNil(t, err)
	assert.EqualError(t, err, "index must not be negative number")
	assert.EqualValues(t, []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}, result)
}

func TestPullAtWithNoDataToPull(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := PullAt(data)
	// ===> []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}, result)
}

func TestRemove(t *testing.T) {
	data := []string{"aa", "bb", "ac", "ad", "ba", "cb", "ac", "vd", "sa", "bb"}
	result, removed, err := Remove(data, func(each string, i int) bool {
		return strings.Contains(each, "a")
	})

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"bb", "cb", "vd", "bb"}, result)
	assert.EqualValues(t, []string{"aa", "ac", "ad", "ba", "ac", "sa"}, removed)
}

func TestRemoveWithInvalidCallback(t *testing.T) {
	data := []string{"aa", "bb", "ac", "ad", "ba", "cb", "ac", "vd", "sa", "bb"}
	result, removed, err := Remove(data, nil)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "callback should be function")
	assert.EqualValues(t, []string{"aa", "bb", "ac", "ad", "ba", "cb", "ac", "vd", "sa", "bb"}, result)
	assert.EqualValues(t, []string{}, removed)
}

func TestReverse(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Reverse(data)
	// ===> []string{"cassandra", "grayson", "damian"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"cassandra", "grayson", "damian"}, result)
}

func TestReverseWithEmptyData(t *testing.T) {
	data := []string{}
	result, err := Reverse(data)
	// ===> []string{}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{}, result)
}

func TestTail(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Tail(data)
	// ===> []string{"grayson", "cassandra"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"grayson", "cassandra"}, result)
}

func TestTake(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Take(data, 2)
	// ===> []string{"damian", "grayson"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "grayson"}, result)
}

func TestTakeRight(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := TakeRight(data, 2)
	// ===> []string{"damian", "grayson"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"grayson", "cassandra"}, result)
}

func TestUnion(t *testing.T) {
	data := []string{"damian", "grayson", "grayson", "cassandra"}
	union1 := []string{"tim", "grayson", "jason", "stephanie"}
	union2 := []string{"monyo"}
	result, err := Union(data, union1, union2)
	// ===> []string{"damian", "grayson", "cassandra", "tim", "jason", "stephanie", "monyo"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "grayson", "cassandra", "tim", "jason", "stephanie", "monyo"}, result)
}

func TestUnionDifferentDataType(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	union1 := []int{1, 2, 3, 4}
	result, err := Union(data, union1)
	_ = result

	assert.NotNil(t, err)
	assert.EqualError(t, err, "data type of each elements between slice must be same")
}

func TestUniq(t *testing.T) {
	data := []string{"damian", "grayson", "grayson", "cassandra"}
	result, err := Uniq(data)
	// ===> []string{"damian", "grayson", "cassandra"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "grayson", "cassandra"}, result)
}

func TestWithout(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}
	result, err := Without(data, "tim", "grayson", "cassandra")
	// ===> []string{"damian", "jason"}

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"damian", "jason"}, result)
}

func TestXor(t *testing.T) {
	t.Skip()

	data := []string{"damian", "grayson", "tim"}
	xorWith := []string{"tim", "jason"}
	result, err := Xor(data, xorWith)
	// ===> []string{"damian", "grayson", "jason"}

	assert.Nil(t, err)
	assert.Equal(t, true, reflect.DeepEqual([]string{"damian", "grayson", "jason"}, result))
}

func TestXorMultipleData(t *testing.T) {
	t.Skip()

	result, err := Xor(
		[]string{"damian", "grayson", "tim", "jason"},
		[]string{"tim", "jason"},
		[]string{"jason"},
	)
	// ===> []string{"damian", "grayson"}
	t.Logf("%#v \n", result)

	assert.Nil(t, err)
	assert.Equal(t, true, reflect.DeepEqual([]string{"damian", "grayson"}, result))
}
