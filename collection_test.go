package gubrak

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountSlice(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Count(data)

	assert.Nil(t, err)
	assert.EqualValues(t, 3, result)
}

func TestCountSliceWithPredicate(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Count(data, func(each string) bool {
		return strings.Contains(each, "d")
	})

	assert.Nil(t, err)
	assert.EqualValues(t, 2, result)
}

func TestCountSliceWithPredicate2(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Count(data, func(each string, i int) bool {
		return len(each) > 6 && i > 1
	})

	assert.Nil(t, err)
	assert.EqualValues(t, 1, result)
}

func TestCountMap(t *testing.T) {
	data := map[string]interface{}{
		"name":   "jason",
		"age":    12,
		"isMale": true,
	}
	result, err := Count(data)

	assert.Nil(t, err)
	assert.EqualValues(t, 3, result)
}

func TestCountMapWithPredicate1(t *testing.T) {
	data := map[string]interface{}{
		"name":   "jason",
		"age":    12,
		"isMale": true,
	}
	result, err := Count(data, func(val interface{}, key string) bool {
		return strings.Contains(strings.ToLower(key), "m")
	})

	assert.Nil(t, err)
	assert.EqualValues(t, 2, result)
}

func TestEachSlice(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	separator := ","
	joinedString := ""

	err := Each(data, func(each string) {
		if joinedString == "" {
			joinedString = each
		} else {
			joinedString = joinedString + separator + each
		}
	})

	assert.Nil(t, err)
	assert.EqualValues(t, "damian,grayson,cassandra", joinedString)
}

func TestEachSliceWithIndex(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	separator := ","
	joinedString := ""

	err := Each(data, func(each string, i int) {
		if i == 0 {
			joinedString = each
		} else {
			joinedString = joinedString + separator + each
		}
	})

	assert.Nil(t, err)
	assert.EqualValues(t, "damian,grayson,cassandra", joinedString)
}

func TestEachSliceStoppable(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra", "tim", "jason", "stephanie"}
	take := 4
	separator := ","
	joinedString := ""

	err := Each(data, func(each string, i int) bool {
		if i >= take {
			return false
		}

		if i == 0 {
			joinedString = each
		} else {
			joinedString = joinedString + separator + each
		}

		return true
	})

	assert.Nil(t, err)
	assert.EqualValues(t, "damian,grayson,cassandra,tim", joinedString)
}

func TestEachSliceWrongLoopParamType(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}

	err := Each(data, func(each map[string]interface{}) {
		// do something
	})

	assert.NotNil(t, err)
	assert.EqualError(t, err, "callback 1st parameter's data type should be same with slice element data type")
}

func TestEachSliceSliceStructData(t *testing.T) {
	type Sample struct {
		Name string
		Age  int
	}

	data := []Sample{
		{Name: "damian", Age: 12},
		{Name: "grayson", Age: 10},
		{Name: "cassandra", Age: 11},
	}
	separator := ","
	joinedString := ""

	err := Each(data, func(each Sample, i int) {
		if i == 0 {
			joinedString = each.Name
		} else {
			joinedString = joinedString + separator + each.Name
		}
	})

	assert.Nil(t, err)
	assert.EqualValues(t, "damian,grayson,cassandra", joinedString)
}

func TestEachCollection(t *testing.T) {
	data := map[string]interface{}{
		"name":   "damian",
		"age":    17,
		"gender": "male",
	}
	separator := ","
	joinedString := ""

	err := Each(data, func(value interface{}, key string) {
		each := fmt.Sprintf("%s: %v", key, value)
		if joinedString == "" {
			joinedString = each
		} else {
			joinedString = joinedString + separator + each
		}
	})

	assert.Nil(t, err)

	for _, each := range strings.Split(joinedString, separator) {
		switch each {
		case "name: damian":
		case "age: 17":
		case "gender: male":
			break
		default:
			assert.Fail(t, "Each function is buggy if used to loop map data")
		}
	}
}

func TestEachCollectionWithCallbackKeyOnly(t *testing.T) {
	data := map[string]interface{}{
		"name":   "damian",
		"age":    17,
		"gender": "male",
	}

	err := Each(data, func(value interface{}, key string) {
		switch key {
		case "name":
		case "age":
		case "gender":
			break
		default:
			assert.Fail(t, "Each function is buggy if used to loop map data")
		}
	})

	assert.Nil(t, err)
}

func TestEachRightSlice(t *testing.T) {
	data := []string{"damian", "grayson", "cassandra"}
	separator := ","
	joinedString := ""

	err := EachRight(data, func(each string) {
		if joinedString == "" {
			joinedString = each
		} else {
			joinedString = joinedString + separator + each
		}
	})

	assert.Nil(t, err)
	assert.EqualValues(t, "cassandra,grayson,damian", joinedString)
}

func TestFilterSlice(t *testing.T) {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := Filter(data, func(each Sample) bool {
		return each.DailyDownloads > 11000
	})

	err = Each(result.([]Sample), func(each Sample) {
		if each.EbookName == "clean code" {
			t.Fail()
		}
	})

	assert.Nil(t, err)
	assert.Len(t, result, 2)
}

func TestFilterZeroMatch(t *testing.T) {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := Filter(data, func(each Sample) bool {
		return each.DailyDownloads > 15000
	})

	assert.Nil(t, err)
	assert.Len(t, result, 0)
}

func TestFilterInvalidLoopEach(t *testing.T) {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	_, err := Filter(data, func(each Sample) int {
		return each.DailyDownloads
	})

	assert.NotNil(t, err)
	assert.EqualError(t, err, "callback return value should be one variable with bool type")
}

func TestFilterCollection(t *testing.T) {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := map[string]int{
		"clean code":       10000,
		"rework":           12000,
		"detective comics": 11500,
	}

	result, err := Filter(data, func(value int, key string) bool {
		return value > 11000
	})

	err = Each(result.(map[string]int), func(value int, key string) {
		if key == "clean code" {
			t.Fail()
		}
	})

	assert.Nil(t, err)
	assert.Len(t, result, 2)
}

func TestFind(t *testing.T) {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	row, err := Find(data, func(each Sample) bool {
		return each.EbookName == "rework"
	})

	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.EqualValues(t, 12000, row.(Sample).DailyDownloads)
}

func TestFindWrongClause(t *testing.T) {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	row, err := Find(data, func(each Sample, i int) bool {
		return each.EbookName == "red hood and the outlaws"
	})

	assert.Nil(t, err)
	assert.Nil(t, row)
}

func TestFindWithFromIndex(t *testing.T) {
	data := []string{"clean code", "rework", "detective comics"}

	row, err := Find(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 1)

	assert.Nil(t, err)
	assert.EqualValues(t, "detective comics", row)
}

func TestFindLast(t *testing.T) {
	data := []string{"clean code", "rework", "detective comics"}

	row, err := FindLast(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	})

	assert.Nil(t, err)
	assert.EqualValues(t, "detective comics", row)
}

func TestFindLastWithFromIndex0(t *testing.T) {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	row, err := FindLast(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 0)

	assert.Nil(t, err)
	assert.EqualValues(t, "clean code", row)
}

func TestFindLastWithFromIndex1(t *testing.T) {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	row, err := FindLast(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 1)

	assert.Nil(t, err)
	assert.EqualValues(t, "clean code", row)
}

func TestFindLastWithFromIndex2(t *testing.T) {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	row, err := FindLast(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 2)

	assert.Nil(t, err)
	assert.EqualValues(t, "detective comics", row)
}

func TestFindLastWithFromIndex3(t *testing.T) {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	row, err := FindLast(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 3)

	assert.Nil(t, err)
	assert.EqualValues(t, "coco", row)
}

func TestFindLastWithFromIndex4(t *testing.T) {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	row, err := FindLast(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 4)

	assert.Nil(t, err)
	assert.EqualValues(t, "coco", row)
}

func TestGroupBy(t *testing.T) {
	type Sample struct {
		Ebook    string
		Category string
	}

	data := []Sample{
		{Ebook: "clean code", Category: "productivity"},
		{Ebook: "rework", Category: "productivity"},
		{Ebook: "detective comics", Category: "comics"},
		{Ebook: "injustice 2", Category: "comics"},
		{Ebook: "dragon ball", Category: "manga"},
		{Ebook: "one piece", Category: "manga"},
	}

	result, err := GroupBy(data, func(each Sample) string {
		return each.Category
	})
	resultParsed := result.(map[string][]Sample)

	assert.Nil(t, err)

	for key, val := range resultParsed {
		switch key {
		case "productivity":
			currResult, err := Filter(val, func(each Sample) bool {
				return each.Ebook == "clean code" || each.Ebook == "rework"
			})
			assert.Nil(t, err)
			assert.Equal(t, true, len(currResult.([]Sample)) > 0, "grouped data has invalid items")
		case "comics":
			currResult, err := Filter(val, func(each Sample) bool {
				return each.Ebook == "detective comics" || each.Ebook == "injustice 2"
			})
			assert.Nil(t, err)
			assert.Equal(t, true, len(currResult.([]Sample)) > 0, "grouped data has invalid items")
		case "manga":
			currResult, err := Filter(val, func(each Sample) bool {
				return each.Ebook == "dragon ball" || each.Ebook == "one piece"
			})
			assert.Nil(t, err)
			assert.Equal(t, true, len(currResult.([]Sample)) > 0, "grouped data has invalid items")
		default:
			assert.Fail(t, "grouped data has different key compared to expected value")
		}
	}
}

func TestGroupByWithFlatDataInt(t *testing.T) {
	data := []int{1, 2, 3, 5, 6, 4, 2, 5, 2}

	result, err := GroupBy(data, func(each int) int {
		return each
	})
	resultParsed := result.(map[int][]int)

	assert.Nil(t, err)

	for key, val := range resultParsed {
		switch key {
		case 1:
			assert.EqualValues(t, []int{1}, val)
		case 2:
			assert.EqualValues(t, []int{2, 2, 2}, val)
		case 3:
			assert.EqualValues(t, []int{3}, val)
		case 5:
			assert.EqualValues(t, []int{5, 5}, val)
		case 6:
			assert.EqualValues(t, []int{6}, val)
		case 4:
			assert.EqualValues(t, []int{4}, val)
		default:
			assert.Fail(t, "grouped data has different key compared to expected value")
		}
	}
}

func TestIncludesSliceString(t *testing.T) {
	data := []string{"damian", "tim", "jason", "grayson"}

	result, err := Includes(data, "tim")
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestIncludesSliceStringWithStartIndex(t *testing.T) {
	data := []string{"damian", "tim", "jason", "grayson"}

	result, err := Includes(data, "tim", 2)
	assert.Nil(t, err)
	assert.False(t, result)
}

func TestIncludesSliceStringWrongSearch(t *testing.T) {
	data := []string{"damian", "tim", "jason", "grayson"}

	result, err := Includes(data, "cassandra")
	assert.Nil(t, err)
	assert.False(t, result)
}

func TestIncludesInterface(t *testing.T) {
	var err error
	var result bool

	data := []interface{}{"name", 12, true}

	result, err = Includes(data, "name")
	assert.Nil(t, err)
	assert.True(t, result)

	result, err = Includes(data, 12)
	assert.Nil(t, err)
	assert.True(t, result)

	result, err = Includes(data, true)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestIncludesString(t *testing.T) {
	data := "damian"

	result, err := Includes(data, "an")
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestIncludesMap(t *testing.T) {
	data := map[string]string{
		"name":  "grayson",
		"hobby": "helping people",
	}

	result, err := Includes(data, "grayson")
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestIncludesMapWrongSearch(t *testing.T) {
	data := map[string]string{
		"name":  "grayson",
		"hobby": "helping people",
	}

	result, err := Includes(data, "batmobile")
	assert.Nil(t, err)
	assert.False(t, result)
}

func TestKeyBy(t *testing.T) {
	type HashMap map[string]string

	data := []HashMap{
		{"name": "grayson", "hobby": "helping people"},
		{"name": "jason", "hobby": "punching people"},
		{"name": "tim", "hobby": "stay awake all the time"},
		{"name": "damian", "hobby": "getting angry"},
	}

	result, err := KeyBy(data, func(each HashMap) string {
		return each["name"]
	})
	assert.Nil(t, err)

	for key, value := range result.(map[string]HashMap) {
		switch key {
		case "grayson":
			assert.EqualValues(t, HashMap{"name": "grayson", "hobby": "helping people"}, value)
		case "jason":
			assert.EqualValues(t, HashMap{"name": "jason", "hobby": "punching people"}, value)
		case "tim":
			assert.EqualValues(t, HashMap{"name": "tim", "hobby": "stay awake all the time"}, value)
		case "damian":
			assert.EqualValues(t, HashMap{"name": "damian", "hobby": "getting angry"}, value)
		default:
			assert.Error(t, errors.New("KeyBy return value is wrong"))
		}
	}
}

func TestMap(t *testing.T) {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	newData, err := Map(data, func(each Sample, i int) string {
		return each.EbookName
	})

	assert.Nil(t, err)
	assert.EqualValues(t, []string{"clean code", "rework", "detective comics"}, newData)
}

func TestMapToNewStruct(t *testing.T) {
	type SampleOne struct {
		EbookName      string
		DailyDownloads int
		IsActive       bool
	}

	type SampleTwo struct {
		Ebook                string
		DownloadsInThousands float32
	}

	data := []SampleOne{
		{EbookName: "clean code", DailyDownloads: 10000, IsActive: true},
		{EbookName: "rework", DailyDownloads: 12000, IsActive: false},
		{EbookName: "detective comics", DailyDownloads: 11500, IsActive: true},
	}

	newData, err := Map(data, func(each SampleOne, i int) SampleTwo {
		ebook := each.EbookName
		if !each.IsActive {
			ebook = fmt.Sprintf("%s (inactive)", each.EbookName)
		}

		return SampleTwo{
			Ebook:                ebook,
			DownloadsInThousands: float32(each.DailyDownloads) / float32(1000),
		}
	})

	assert.Nil(t, err)
	assert.EqualValues(t, []SampleTwo{
		{Ebook: "clean code", DownloadsInThousands: 10},
		{Ebook: "rework (inactive)", DownloadsInThousands: 12},
		{Ebook: "detective comics", DownloadsInThousands: 11.5},
	}, newData)
}

func TestOrderByString(t *testing.T) {
	type HashMap map[string]string

	data := []HashMap{
		{"name": "tim", "hobby": "stay awake all the time"},
		{"name": "grayson", "hobby": "helping people"},
		{"name": "damian", "hobby": "getting angry"},
		{"name": "jason", "hobby": "punching people"},
	}

	result, err := OrderBy(data, func(each HashMap) string {
		return each["name"]
	})
	resultParsed := result.([]HashMap)

	for i, each := range resultParsed {
		switch i {
		case 0:
			assert.EqualValues(t, "damian", each["name"])
		case 1:
			assert.EqualValues(t, "grayson", each["name"])
		case 2:
			assert.EqualValues(t, "jason", each["name"])
		case 3:
			assert.EqualValues(t, "tim", each["name"])
		}
	}

	assert.Nil(t, err)
}

func TestOrderByNumber(t *testing.T) {
	type HashMap map[string]interface{}

	data := []HashMap{
		{"name": "tim", "hobby": "stay awake all the time", "age": 20},
		{"name": "grayson", "hobby": "helping people", "age": 24},
		{"name": "damian", "hobby": "getting angry", "age": 17},
		{"name": "jason", "hobby": "punching people", "age": 22},
	}
	result, err := OrderBy(data, func(each HashMap) int {
		return each["age"].(int)
	})

	resultParsed := result.([]HashMap)

	for i, each := range resultParsed {
		switch i {
		case 0:
			assert.EqualValues(t, "damian", each["name"])
		case 1:
			assert.EqualValues(t, "tim", each["name"])
		case 2:
			assert.EqualValues(t, "jason", each["name"])
		case 3:
			assert.EqualValues(t, "grayson", each["name"])
		}
	}

	assert.Nil(t, err)
}

func TestOrderDescending(t *testing.T) {
	type HashMap map[string]interface{}

	data := []HashMap{
		{"name": "tim", "hobby": "stay awake all the time", "age": 20},
		{"name": "grayson", "hobby": "helping people", "age": 24},
		{"name": "damian", "hobby": "getting angry", "age": 17},
		{"name": "jason", "hobby": "punching people", "age": 22},
	}
	result, err := OrderBy(data, func(each HashMap) int {
		return each["age"].(int)
	}, false)

	resultParsed := result.([]HashMap)

	for i, each := range resultParsed {
		switch i {
		case 0:
			assert.EqualValues(t, "grayson", each["name"])
		case 1:
			assert.EqualValues(t, "jason", each["name"])
		case 2:
			assert.EqualValues(t, "tim", each["name"])
		case 3:
			assert.EqualValues(t, "damian", each["name"])
		}
	}

	assert.Nil(t, err)
}

func TestOrderNotAsyncSort(t *testing.T) {
	type HashMap map[string]interface{}

	data := []HashMap{
		{"name": "tim", "hobby": "stay awake all the time", "age": 20},
		{"name": "grayson", "hobby": "helping people", "age": 24},
		{"name": "damian", "hobby": "getting angry", "age": 17},
		{"name": "jason", "hobby": "punching people", "age": 22},
	}
	result, err := OrderBy(data, func(each HashMap) int {
		return each["age"].(int)
	}, true, false)

	resultParsed := result.([]HashMap)

	for i, each := range resultParsed {
		switch i {
		case 0:
			assert.EqualValues(t, "damian", each["name"])
		case 1:
			assert.EqualValues(t, "tim", each["name"])
		case 2:
			assert.EqualValues(t, "jason", each["name"])
		case 3:
			assert.EqualValues(t, "grayson", each["name"])
		}
	}

	assert.Nil(t, err)
}

func TestPartition(t *testing.T) {
	type HashMap map[string]interface{}

	data := []HashMap{
		{"name": "grayson", "isMale": true},
		{"name": "jason", "isMale": true},
		{"name": "barbara", "isMale": false},
		{"name": "tim", "isMale": true},
		{"name": "cassandra", "isMale": false},
		{"name": "stephanie", "isMale": false},
		{"name": "damian", "isMale": true},
		{"name": "duke", "isMale": true},
	}

	resultTruthy, resultFalsey, err := Partition(data, func(each HashMap) bool {
		return each["isMale"].(bool)
	})

	assert.Nil(t, err)

	assert.EqualValues(t, []HashMap{
		{"name": "grayson", "isMale": true},
		{"name": "jason", "isMale": true},
		{"name": "tim", "isMale": true},
		{"name": "damian", "isMale": true},
		{"name": "duke", "isMale": true},
	}, resultTruthy)

	assert.EqualValues(t, []HashMap{
		{"name": "barbara", "isMale": false},
		{"name": "cassandra", "isMale": false},
		{"name": "stephanie", "isMale": false},
	}, resultFalsey)
}

func TestReduceSliceNumber(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result, err := Reduce(data, func(current, each int) int {
		return current + each
	}, 0)

	assert.Nil(t, err)
	assert.EqualValues(t, 55, result)
}

func TestReduceSlice(t *testing.T) {
	type HashMap map[string]interface{}

	data := [][]interface{}{
		{"name", "grayson"},
		{"age", 21},
		{"isMale", true},
	}
	result, err := Reduce(data, func(current HashMap, each []interface{}, i int) HashMap {
		current[each[0].(string)] = each[1]
		return current
	}, HashMap{})

	assert.Nil(t, err)
	assert.EqualValues(t, HashMap{
		"name":   "grayson",
		"age":    21,
		"isMale": true,
	}, result)
}

func TestReduceCollection(t *testing.T) {
	type HashMap map[string]interface{}

	data := HashMap{
		"name":   "grayson",
		"age":    21,
		"isMale": true,
	}

	result, err := Reduce(data, func(current string, value interface{}, key string) string {
		if current == "" {
			current = fmt.Sprintf("%s: %v", key, value)
		} else {
			current = fmt.Sprintf("%s, %s: %v", current, key, value)
		}

		return current
	}, "")

	assert.Nil(t, err)
	assert.True(t,
		assert.Contains(t, result, "name: grayson"),
		assert.Contains(t, result, "age: 21"),
		assert.Contains(t, result, "isMale: true"),
	)
}

func TestReject(t *testing.T) {
	type Book struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Book{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := Reject(data, func(each Book) bool {
		return each.DailyDownloads > 11000
	})

	Each(result.([]Book), func(each Book) {
		if each.EbookName == "rework" || each.EbookName == "detective comics" {
			t.Fail()
		}
	})

	assert.Nil(t, err)
	assert.Len(t, result, 1)
}

func TestSample(t *testing.T) {
	type Book struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Book{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := Sample(data)
	resultParsed := result.(Book)

	assert.Nil(t, err)
	switch resultParsed.EbookName {
	case "clean code":
	case "rework":
	case "detective comics":
		break
	default:
		t.Fail()
	}
}

func TestSampleSize(t *testing.T) {
	type Book struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Book{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := SampleSize(data, 2)
	resultParsed := result.([]Book)

	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(resultParsed))

	for _, each := range resultParsed {
		switch each.EbookName {
		case "clean code":
		case "rework":
		case "detective comics":
			break
		default:
			t.Fail()
		}
	}
}

func TestShuffle(t *testing.T) {
	type Book struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Book{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := Shuffle(data)
	resultParsed := result.([]Book)

	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(resultParsed))

	for _, each := range resultParsed {
		switch each.EbookName {
		case "clean code":
		case "rework":
		case "detective comics":
			break
		default:
			t.Fail()
		}
	}
}

func TestSizeSlice(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	result, err := Size(data)

	assert.Nil(t, err)
	assert.EqualValues(t, len(data), result)
}

func TestSizeSliceString(t *testing.T) {
	data := "bruce"
	result, err := Size(data)

	assert.Nil(t, err)
	assert.EqualValues(t, len(data), result)
}

func TestSizeCollection(t *testing.T) {
	data := map[string]interface{}{
		"name":   "noval",
		"age":    24,
		"isMale": true,
	}
	result, err := Size(data)

	assert.Nil(t, err)
	assert.EqualValues(t, len(data), result)
}
