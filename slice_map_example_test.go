package gubrak

import (
	"fmt"
	"log"
	"strings"
)

func ExampleChunk_chunk1() {
	data := []int{1, 2, 3, 4, 5}
	size := 2

	result, err := Chunk(data, size)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> [][]int{ { 1, 2 }, { 3, 4 }, { 5 } }
}

func ExampleChunk_chunk2() {
	data := []string{"a", "b", "c", "d", "e"}
	size := 3

	result, err := Chunk(data, size)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> [][]string{ { "a", "b", "c" }, { "d", "e" } }
}

func ExampleChunk_chunk3() {
	data := []interface{}{
		3.2, "a", -1,
		make([]byte, 0),
		map[string]int{"b": 2},
		[]string{"a", "b", "c"},
	}
	size := 3

	result, err := Chunk(data, size)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		  [][]interface{}{
			{ 3.2, "a" },
			{ -1, []uint8{} },
			{ map[string]int{ "b":2 }, []string{ "a", "b", "c" } },
		  }
	*/
}

func ExampleCompact_compact1() {
	data := []int{-2, -1, 0, 1, 2}

	result, err := Compact(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ -2, -1, 1, 2 }
}

func ExampleCompact_compact2() {
	data := []string{"a", "b", "", "d"}

	result, err := Compact(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "a", "b", "d" }
}

func ExampleCompact_compact3() {
	data := []interface{}{-2, 0, 1, 2, false, true, "", "hello", nil}

	result, err := Compact(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []interface{}{ -2, 1, 2, true, "hello" }
}

func ExampleCompact_compact4() {
	item1, item2, item3 := "a", "b", "c"
	data := []*string{&item1, nil, &item2, nil, &item3}

	result, err := Compact(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []*string{ (*string)(0xc42000e1e0), (*string)(0xc42000e1f0), (*string)(0xc42000e200) }
}

func ExampleConcat_concat1() {
	data := []int{1, 2, 3, 4}
	dataConcat1 := []int{4, 6, 7}
	dataConcat2 := []int{8, 9}

	result, err := Concat(data, dataConcat1, dataConcat2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 5, 6, 7, 8, 9 }
}

func ExampleConcat_concat2() {
	data := []string{"my"}
	dataConcat1 := []string{"name", "is"}
	dataConcat2 := []string{"jason", "todd"}

	result, err := Concat(data, dataConcat1, dataConcat2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "my", "name", "is", "jason", "todd" }
}

func ExampleCount_countMap1() {
	data := map[string]interface{}{
		"name":   "jason",
		"age":    12,
		"isMale": true,
	}
	result, err := Count(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 3
}

func ExampleCount_countMap2() {
	data := map[string]interface{}{
		"name":   "jason",
		"age":    12,
		"isMale": true,
	}

	result, err := Count(data, func(val interface{}, key string) bool {
		return strings.Contains(strings.ToLower(key), "m")
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 2
}

func ExampleCount_countMap3() {
	data := map[string]interface{}{
		"name":   "jason",
		"age":    12,
		"isMale": true,
	}

	result, err := Count(data, func(val interface{}, key string, i int) bool {
		return strings.Contains(strings.ToLower(key), "m") && i > 1
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 1
}

func ExampleCount_countSlice1() {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Count(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 3
}

func ExampleCount_countSlice2() {
	data := []string{"damian", "grayson", "cassandra"}

	result, err := Count(data, func(each string) bool {
		return strings.Contains(each, "d")
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 2
}

func ExampleCount_countSlice3() {
	data := []string{"damian", "grayson", "cassandra"}

	result, err := Count(data, func(each string, i int) bool {
		return len(each) > 6 && i > 1
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 1
}

func ExampleDifference_difference1() {
	data := []int{1, 2, 3, 4, 4, 6, 7}
	dataDiff := []int{2, 7}

	result, err := Difference(data, dataDiff)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 3, 4, 4, 6 }
}

func ExampleDifference_difference2() {
	data := []string{"a", "b", "b", "c", "d", "e", "f", "g", "h"}
	dataDiff1 := []string{"b", "d"}
	dataDiff2 := []string{"e", "f", "h"}

	result, err := Difference(data, dataDiff1, dataDiff2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "a", "c", "g" }
}

func ExampleDifference_difference3() {
	data := []float64{1.1, 1.11, 1.2, 2.3, 3.0, 3, 4.0, 4.00000, 4.000000001}
	dataDiff1 := []float64{1.1, 3}
	dataDiff2 := []float64{4.000000001}

	result, err := Difference(data, dataDiff1, dataDiff2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.11, 1.2, 2.3, 4, 4 }
}

func ExampleDrop_drop1() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	n := 1

	result, err := Drop(data, n)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 2, 3, 4, 4, 5, 6 }
}

func ExampleDrop_drop2() {
	data := []string{"a", "b", "c", "d", "e", "f"}
	n := 3

	result, err := Drop(data, n)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "d", "e", "f" }
}

func ExampleDropRight_dropRight1() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	n := 1

	result, err := DropRight(data, n)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 4, 5 }
}

func ExampleDropRight_dropRight2() {
	data := []string{"a", "b", "c", "d", "e", "f"}
	n := 3

	result, err := DropRight(data, n)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "a", "b", "c" }
}

func ExampleEach_eachMap1() {
	data := map[string]interface{}{
		"name":   "damian",
		"age":    17,
		"gender": "male",
	}

	err := Each(data, func(value interface{}, key string) {
		fmt.Printf("%s: %v \n", key, value)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEach_eachMap2() {
	data := map[string]interface{}{
		"name":   "damian",
		"age":    17,
		"gender": "male",
	}

	err := Each(data, func(value interface{}, key string, i int) {
		fmt.Printf("key: %s, value: %v, index: %d \n", key, value, i)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEach_eachSlice1() {
	data := []string{"damian", "grayson", "cassandra"}

	err := Each(data, func(each string) {
		fmt.Println(each)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEach_eachSlice2() {
	data := []string{"damian", "grayson", "cassandra"}

	err := Each(data, func(each string, i int) {
		fmt.Printf("element %d: %s \n", i, each)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEach_eachSlice3() {
	type Sample struct {
		Name string
		Age  int
	}

	data := []Sample{
		{Name: "damian", Age: 12},
		{Name: "grayson", Age: 10},
		{Name: "cassandra", Age: 11},
	}

	err := Each(data, func(each Sample) {
		fmt.Printf("name: %s, age: %d \n", each.Name, each.Age)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEach_eachSlice4() {
	data := []string{"damian", "grayson", "cassandra", "tim", "jason", "stephanie"}

	err := Each(data, func(each string, i int) bool {
		if i > 3 { // will stop after fourth loop
			return false
		}

		fmt.Println(each)
		return true
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEachRight_eachRightMap1() {
	data := map[string]interface{}{
		"name":   "damian",
		"age":    17,
		"gender": "male",
	}

	err := EachRight(data, func(value interface{}, key string) {
		fmt.Printf("%s: %v \n", key, value)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEachRight_eachRightMap2() {
	data := map[string]interface{}{
		"name":   "damian",
		"age":    17,
		"gender": "male",
	}

	err := EachRight(data, func(value interface{}, key string, i int) {
		fmt.Printf("key: %s, value: %v, index: %d \n", key, value, i)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEachRight_eachRightSlice1() {
	data := []string{"damian", "grayson", "cassandra"}

	err := EachRight(data, func(each string) {
		fmt.Println(each)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEachRight_eachRightSlice2() {
	data := []string{"damian", "grayson", "cassandra"}

	err := EachRight(data, func(each string, i int) {
		fmt.Printf("element %d: %s \n", i, each)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEachRight_eachRightSlice3() {
	type Sample struct {
		Name string
		Age  int
	}

	data := []Sample{
		{Name: "damian", Age: 12},
		{Name: "grayson", Age: 10},
		{Name: "cassandra", Age: 11},
	}

	err := EachRight(data, func(each Sample) {
		fmt.Printf("name: %s, age: %d \n", each.Name, each.Age)
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEachRight_eachRightSlice4() {
	data := []string{"damian", "grayson", "cassandra", "tim", "jason", "stephanie"}

	err := EachRight(data, func(each string, i int) bool {
		if i > 3 { // will stop after fourth loop
			return false
		}

		fmt.Println(each)
		return true
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleFill_fill1() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	replacement := 9

	result, err := Fill(data, replacement)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 9, 9, 9, 9, 9, 9, 9 }
}

func ExampleFill_fill2() {
	data := []string{"grayson", "jason", "tim", "damian"}
	replacement := "alfred"
	start := 2

	result, err := Fill(data, replacement, start)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ "grayson", "jason", "alfred", "alfred" }
}

func ExampleFill_fill3() {
	data := []float64{1, 2.2, 3.0002, 4, 4, 5.12, 6}
	replacement := float64(9)
	start, end := 3, 5

	result, err := Fill(data, replacement, start, end)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1, 2.2, 3.0002, 9, 9, 5.12, 6 }
}

func ExampleFilter_filterMap() {
	data := map[string]int{
		"clean code":       10000,
		"rework":           12000,
		"detective comics": 11500,
	}

	result, err := Filter(data, func(value int, key string) bool {
		return value > 11000
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		  map[string]int{
			"rework":           12000,
			"detective comics": 11500,
		  }
	*/
}

func ExampleFilter_filterSlice() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		  []Sample{
			{ EbookName: "rework", DailyDownloads: 12000 },
			{ EbookName: "detective comics", DailyDownloads: 11500 },
		  }
	*/
}

func ExampleFind_find1() {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := Find(data, func(each Sample) bool {
		return each.EbookName == "rework"
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> Sample { EbookName: "rework", DailyDownloads: 12000 }
}

func ExampleFind_find2() {
	data := []string{"clean code", "rework", "detective comics"}

	result, err := Find(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 1)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "detective comics"
}

func ExampleFindIndex_findIndex1() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}
	predicate := func(each string) bool {
		return each == "tim"
	}

	result, err := FindIndex(data, predicate)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 3
}

func ExampleFindIndex_findIndex2() {
	data := []int{-2, -1, 0, 1, 2}

	result, err := FindIndex(data, func(each int) bool {
		return each == 4
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> -1
}

func ExampleFindIndex_findIndex3() {
	data := []float64{1, 1.1, 1.2, 1.200001, 1.2000000001, 1.3}

	result, err := FindIndex(data, func(each float64) bool {
		return each == 1.2000000001
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 4
}

func ExampleFindIndex_findIndex4() {
	data := []int{1, 2, 3, 3, 4, 5}
	predicate := func(each int) bool {
		return each == 3
	}
	fromIndex := 2

	result, err := FindIndex(data, predicate, fromIndex)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 2
}

func ExampleFindIndex_findIndex5() {
	data := []int{1, 2, 3, 3, 4, 5}
	predicate := func(each int) bool {
		return each == 3
	}
	fromIndex := 3

	result, err := FindIndex(data, predicate, fromIndex)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 3
}

func ExampleFindLast_findLast1() {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := FindLast(data, func(each Sample) bool {
		return strings.Contains(each.EbookName, "co")
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> Sample { EbookName: "detective comics", DailyDownloads: 11500 }
}

func ExampleFindLast_findLast2() {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	result, err := FindLast(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "detective comics"
}

func ExampleFindLast_findLast3() {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	result, err := FindLast(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 3)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "coco"
}

func ExampleFindLastIndex_findLastIndex1() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}

	result, err := FindLastIndex(data, func(each string) bool {
		return each == "tim"
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 4
}

func ExampleFindLastIndex_findLastIndex2() {
	data := []int{1, 2, 2, 3, 3, 4, 5}
	predicate := func(each int) bool {
		return each == 3
	}
	fromIndex := 4

	result, err := FindLastIndex(data, predicate, fromIndex)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 4
}

func ExampleFindLastIndex_findLastIndex3() {
	data := []int{1, 2, 3, 3, 4, 5}
	predicate := func(each int) bool {
		return each == 3
	}
	fromIndex := 3

	result, err := FindLastIndex(data, predicate, fromIndex)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 3
}

func ExampleFindLastIndex_findLastIndex4() {
	data := []int{1, 2, 3, 3, 4, 5}
	predicate := func(each int) bool {
		return each == 3
	}
	fromIndex := 2

	result, err := FindLastIndex(data, predicate, fromIndex)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> -1
}

func ExampleFirst_first1() {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := First(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "damian"
}

func ExampleFirst_first2() {
	data := []string{}
	result, err := First(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> nil
}

func ExampleFromPairs_fromPairs1() {
	data := []interface{}{
		[]interface{}{"a", 1},
		[]interface{}{"b", 2},
	}

	result, err := FromPairs(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		  map[interface{}]interface{}{
			"a": 1,
			"b": 2,
		  }
	*/
}

func ExampleFromPairs_fromPairs2() {
	data := []interface{}{
		[]interface{}{true, []int{1, 2, 3}},
		[]interface{}{false, []string{"damian", "grayson"}},
	}

	result, err := FromPairs(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		map[interface{}]interface{}{
		  true: []int{ 1, 2, 3 },
		  false: []string{ "damian", "grayson" },
		}
	*/
}

func ExampleGroupBy_groupBy1() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		  map[string][]main.Sample {
			"productivity": []main.Sample {
			  { Ebook: "clean code", Category: "productivity" },
			  { Ebook: "rework", Category: "productivity" },
			},
			"comics":       []main.Sample {
			  { Ebook: "detective comics", Category: "comics"},
			  { Ebook: "injustice 2", Category: "comics"},
			},
			"manga":        []main.Sample {
			  { Ebook: "dragon ball", Category: "manga" },
			  { Ebook: "one piece", Category: "manga"},
			},
		  }
	*/
}

func ExampleGroupBy_groupBy2() {
	data := []int{1, 2, 3, 5, 6, 4, 2, 5, 2}

	result, err := GroupBy(data, func(each int) int {
		return each
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
	   map[int][]int{
	     5: []int{ 5, 5 },
	     6: []int{ 6 },
	     4: []int{ 4 },
	     1: []int{ 1 },
	     2: []int{ 2, 2, 2 },
	     3: []int{ 3 },
	   }
	*/
}

func ExampleIncludes_includesMap1() {
	data := map[string]string{
		"name":  "grayson",
		"hobby": "helping people",
	}

	result, err := Includes(data, "grayson")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> true
}

func ExampleIncludes_includesMap2() {
	data := map[string]string{
		"name":  "grayson",
		"hobby": "helping people",
	}

	result, err := Includes(data, "batmobile")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> false
}

func ExampleIncludes_includesSlice1() {
	data := []string{"damian", "tim", "jason", "grayson"}
	result, err := Includes(data, "tim")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> true
}

func ExampleIncludes_includesSlice2() {
	data := []string{"damian", "tim", "jason", "grayson"}
	result, err := Includes(data, "tim", 2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> false
}

func ExampleIncludes_includesSlice3() {
	data := []string{"damian", "tim", "jason", "grayson"}
	result, err := Includes(data, "cassandra")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> false
}

func ExampleIncludes_includesSlice4() {
	data := []interface{}{"name", 12, true}

	Includes(data, "name") // ===> true
	Includes(data, 12)     // ===> true
	Includes(data, true)   // ===> true
}

func ExampleIncludes_includesSlice5() {
	Includes("damian", "an") // ===> true
}

func ExampleIndexOf_indexOf1() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}
	IndexOf(data, "duke")    // ===> -1
	IndexOf(data, "tim")     // ===> 3
	IndexOf(data, "tim", 4)  // ===> 4
	IndexOf(data, "tim", -4) // ===> 3
	IndexOf(data, "tim", -3) // ===> 4
	IndexOf(data, "tim", -2) // ===> -1
}

func ExampleIndexOf_indexOf2() {
	data := []float64{2.1, 2.2, 3, 3.00000, 3.1, 3.9, 3.95}
	IndexOf(data, 2.2)           // ===> 1
	IndexOf(data, 3)             // ===> -1
	IndexOf(data, float64(3))    // ===> 2 (because 3 is detected as int32, not float64)
	IndexOf(data, float64(3), 2) // ===> 2
	IndexOf(data, float64(3), 3) // ===> 3
}

func ExampleIndexOf_indexOf3() {
	data := []interface{}{"jason", 24, true}
	IndexOf(data, 24)     // ===> 1
	IndexOf(data, 24, -1) // ===> -1
}

func ExampleInitial_initial1() {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Initial(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "grayson" }
}

func ExampleInitial_initial2() {
	data := []int{1, 2, 3, 4, 5}
	result, err := Initial(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4 }
}

func ExampleInitial_initial3() {
	data := []map[string]string{{"name": "jason"}}
	result, err := Initial(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []map[string]string{}
}

func ExampleInitial_initial4() {
	data := []float64{}
	result, err := Initial(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{}
}

func ExampleIntersection_intersection1() {
	result, err := Intersection(
		[]string{"damian", "grayson", "cassandra", "tim", "tim", "jason"},
		[]string{"cassandra", "tim", "jason"},
		[]string{"cassandra", "jason"},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "cassandra", "jason" }
}

func ExampleIntersection_intersection2() {
	result, err := Intersection(
		[]float64{0.8, 0.8001, 0.999, 1, 1.0, 1.000001, 1.1000000, 1.1001, 1.2, 1.33, 1.4},
		[]float64{0.8, 0.8001, 0.999, 1, 1.0, 1.000001, 1.1000000, 1.2, 1.33},
		[]float64{1.1000000, 1.2, 0.8001, 0.999, 1.33, 1, 1.0, 1.000001},
		[]float64{1.2, 0.8001, 0.999, 1.33, 1.000092},
		[]float64{0.8001, 0.999, 1.33, 1.400001},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 0.8001, 0.999, 1.33 }
}

func ExampleJoin_join1() {
	data := []string{"damian", "grayson", "cassandra"}
	separator := " - "

	result, err := Join(data, separator)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "damian - grayson - cassandra"
}

func ExampleJoin_join2() {
	data := []int{1, 2, 3, 4}
	separator := ", "

	result, err := Join(data, separator)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "1, 2, 3, 4"
}

func ExampleKeyBy_keyBy() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
	   map[string]main.HashMap {
	     "grayson": main.HashMap{ "hobby": "helping people", "name": "grayson" },
	     "jason":   main.HashMap{ "name": "jason", "hobby": "punching people" },
	     "tim":     main.HashMap{ "name": "tim", "hobby": "stay awake all the time" },
	     "damian":  main.HashMap{ "name": "damian", "hobby": "getting angry" },
	   }
	*/
}

func ExampleLast_last1() {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Last(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "cassandra"
}

func ExampleLast_last2() {
	data := []int{1}
	result, err := Last(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 1
}

func ExampleLast_last3() {
	data := []string{}
	result, err := Last(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> nil
}

func ExampleLastIndexOf_lastIndexOf1() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}

	LastIndexOf(data, "duke")    // ===> -1
	LastIndexOf(data, "tim")     // ===> 4
	LastIndexOf(data, "tim", 4)  // ===> 4
	LastIndexOf(data, "tim", -4) // ===> 3
	LastIndexOf(data, "tim", -3) // ===> 4
	LastIndexOf(data, "tim", -2) // ===> 4
}

func ExampleLastIndexOf_lastIndexOf2() {
	data := []float64{2.1, 2.2, 3, 3.00000, 3.1, 3.9, 3.95}

	LastIndexOf(data, 2.2)           // ===> 1
	LastIndexOf(data, 3)             // ===> -1 (because 3 is detected as int32, not float64)
	LastIndexOf(data, float64(3))    // ===> 3
	LastIndexOf(data, float64(3), 2) // ===> 2
	LastIndexOf(data, float64(3), 3) // ===> 3
}

func ExampleLastIndexOf_lastIndexOf3() {
	data := []interface{}{"jason", 24, true}

	LastIndexOf(data, 24)     // ===> 1
	LastIndexOf(data, 24, -1) // ===> 1
}

func ExampleMap_map1() {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := Map(data, func(each Sample, i int) string {
		return each.EbookName
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "clean code", "rework", "detective comics" }
}

func ExampleMap_map2() {
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

	result, err := Map(data, func(each SampleOne, i int) SampleTwo {
		ebook := each.EbookName
		if !each.IsActive {
			ebook = fmt.Sprintf("%s (inactive)", each.EbookName)
		}

		downloadsInThousands := float32(each.DailyDownloads) / float32(1000)

		return SampleTwo{Ebook: ebook, DownloadsInThousands: downloadsInThousands}
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*

		  []SampleTwo {
			{ Ebook: "clean code", DownloadsInThousands: 10 },
			{ Ebook: "rework (inactive)", DownloadsInThousands: 12 },
			{ Ebook: "detective comics", DownloadsInThousands: 11.5 },
		  }

	*/
}

func ExampleNth_nth1() {
	data := []string{"grayson", "jason", "tim", "damian"}

	Nth(data, 1)  // ===> "jason"
	Nth(data, 2)  // ===> "tim"
	Nth(data, -1) // ===> "damian"
}

func ExampleNth_nth2() {
	data := []int{1, 2, 3, 4, 5}
	result, err := Nth(data, 4)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 5
}

func ExampleOrderBy_orderBy1() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
	   []main.HashMap{
	     { "name": "damian", "hobby": "getting angry" },
	     { "name": "grayson", "hobby": "helping people" },
	     { "name": "jason", "hobby": "punching people" },
	     { "name": "tim", "hobby": "stay awake all the time" },
	   }
	*/
}

func ExampleOrderBy_orderBy2() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		[]main.HashMap{
		  { "age": 17, "hobby": "getting angry", "name": "damian" },
		  { "age": 20, "name": "tim", "hobby": "stay awake all the time" },
		  { "age": 22, "name": "jason", "hobby": "punching people" },
		  { "age": 24, "name": "grayson", "hobby": "helping people" },
		}
	*/
}

func ExampleOrderBy_orderBy3() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
	   []main.HashMap{
	     { "age": 24, "name": "grayson", "hobby": "helping people" },
	     { "age": 22, "name": "jason", "hobby": "punching people" },
	     { "age": 20, "name": "tim", "hobby": "stay awake all the time" },
	     { "age": 17, "name": "damian", "hobby": "getting angry" },
	   }
	*/
}

func ExampleOrderBy_orderBy4() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
	   []main.HashMap{
	     { "age": 17, "name": "damian", "hobby": "getting angry" },
	     { "age": 20, "name": "tim", "hobby": "stay awake all the time" },
	     { "age": 22, "name": "jason", "hobby": "punching people" },
	     { "age": 24, "name": "grayson", "hobby": "helping people" },
	   }
	*/
}

func ExamplePartition_partition() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%#v \n", resultTruthy)
	/*
	   []HashMap {
	     { "name": "grayson", "isMale": true },
	     { "name": "jason", "isMale": true },
	     { "name": "tim", "isMale": true },
	     { "name": "damian", "isMale": true },
	     { "name": "duke", "isMale": true },
	   }
	*/

	fmt.Printf("%#v \n", resultFalsey)
	/*
	   []HashMap {
	     { "name": "barbara", "isMale": false },
	     { "name": "cassandra", "isMale": false },
	     { "name": "stephanie", "isMale": false },
	   }
	*/
}

func ExamplePull_pull1() {
	data := []int{1, 2, 3, 4, 5, 6}
	result, err := Pull(data, 3)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 4, 5, 6 }
}

func ExamplePull_pull2() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}
	result, err := Pull(data, 2.1, 3.2, 6.3)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 4.2, 5.2 }
}

func ExamplePull_pull3() {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := Pull(data, "grayson", "tim")

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "cassandra", "jason", "stephanie" }
}

func ExamplePullAll_pullAll1() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}
	exclude := []float64{2.1, 3.2, 6.3}

	result, err := PullAll(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 4.2, 5.2 }
}

func ExamplePullAll_pullAll2() {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	exclude := []string{"grayson", "tim"}

	result, err := PullAll(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "cassandra", "jason", "stephanie" }
}

func ExamplePullAt_pullAt() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}

	result, err := PullAt(data, 1, 3)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 3.2, 5.2, 6.3 }
}

func ExampleReduce_reduceMap1() {
	type HashMap map[string]interface{}

	data := HashMap{
		"name":   "grayson",
		"age":    21,
		"isMale": true,
	}

	result, err := Reduce(data, func(accumulator string, value interface{}, key string) string {
		if accumulator == "" {
			accumulator = fmt.Sprintf("%s: %v", key, value)
		} else {
			accumulator = fmt.Sprintf("%s, %s: %v", accumulator, key, value)
		}

		return accumulator
	}, "")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "name: grayson, age: 21, isMale: true"
}

func ExampleReduce_reduceSlice1() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result, err := Reduce(data, func(accumulator, each int) int {
		return accumulator + each
	}, 0)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 55
}

func ExampleReduce_reduceSlice2() {
	type HashMap map[string]interface{}

	data := [][]interface{}{
		{"name", "grayson"},
		{"age", 21},
		{"isMale", true},
	}

	result, err := Reduce(data, func(accumulator HashMap, each []interface{}, i int) HashMap {
		accumulator[each[0].(string)] = each[1]
		return accumulator
	}, HashMap{})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
	   HashMap {
	     "name":   "grayson",
	     "age":    21,
	     "isMale": true,
	   }
	*/
}

func ExampleReject_rejectMap() {
	data := map[string]int{
		"clean code":       10000,
		"rework":           12000,
		"detective comics": 11500,
	}

	result, err := Reject(data, func(value int, key string) bool {
		return value > 11000
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> map[string]int{ "clean code": 10000 }
}

func ExampleReject_rejectSlice() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		  []Book{
			{ EbookName: "clean code", DailyDownloads: 10000 },
		  }
	*/
}

func ExampleRemove_remove1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	result, removed, err := Remove(data, func(each string) bool {
		return strings.Contains(each, "m")
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)  // ===> []string{ "jason", "grayson" }
	fmt.Println(removed) // ===> []string{ "damian", "tim" }
}

func ExampleRemove_remove2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result, removed, err := Remove(data, func(each int) bool {
		return each%2 == 0
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)  // ===> []int{ 1, 3, 5, 7, 9 }
	fmt.Println(removed) // ===> []int{ 2, 4, 6, 8 }
}

func ExampleReverse_reverse1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	result, err := Reverse(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "tim", "grayson", "damian", "jason" }
}

func ExampleReverse_reverse2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result, err := Reverse(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 9, 8, 7, 6, 5, 4, 3, 2, 1 }
}

func ExampleSample_sample() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		  the result can be this:
			{ EbookName: "clean code", DailyDownloads: 10000 },

		  this:
			{ EbookName: "rework", DailyDownloads: 12000 },

		  or this:
			{ EbookName: "detective comics", DailyDownloads: 11500 },
	*/
}

func ExampleSampleSize_sampleSize() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		  the result can be this:
			[]Book{
			  { EbookName: "clean code", DailyDownloads: 10000 },
			  { EbookName: "rework", DailyDownloads: 12000 },
			}

		  this:
			[]Book{
			  { EbookName: "rework", DailyDownloads: 12000 },
			  { EbookName: "detective comics", DailyDownloads: 11500 },
			}

		  or this:
			[]Book{
			  { EbookName: "clean code", DailyDownloads: 10000 },
			  { EbookName: "detective comics", DailyDownloads: 11500 },
			}
	*/
}

func ExampleShuffle_shuffle1() {
	data := []int{1, 2, 3, 4}
	result, err := Shuffle(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		the result can be this:
		  []int{ 1, 4, 2, 3 }

		this:
		  []int{ 4, 1, 2, 3 }

		or this:
		  []int{ 4, 1, 3, 2 }

		or this:
		  []int{ 3, 4, 1, 2 }

		or ... any other possibilities.
	*/
}

func ExampleShuffle_shuffle2() {
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
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	/*
		  the result can be this:
			[]Book {
			  { EbookName: "detective comics", DailyDownloads: 11500 },
			  { EbookName: "clean code", DailyDownloads: 10000 },
			  { EbookName: "rework", DailyDownloads: 12000 },
			}

		  this:
			[]Book {
			  { EbookName: "clean code", DailyDownloads: 10000 },
			  { EbookName: "detective comics", DailyDownloads: 11500 },
			  { EbookName: "rework", DailyDownloads: 12000 },
			}

		  or this:
			[]Book {
			  { EbookName: "rework", DailyDownloads: 12000 },
			  { EbookName: "detective comics", DailyDownloads: 11500 },
			  { EbookName: "clean code", DailyDownloads: 10000 },
			}
	*/
}

func ExampleSize_sizeSlice() {
	Size([]int{1, 2, 3, 4, 5}) // ===> 5
	Size("bruce")              // ===> 5
}

func ExampleSize_sizeMap() {
	data := map[string]interface{}{
		"name":   "noval",
		"age":    24,
		"isMale": true,
	}

	result, err := Size(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 3
}

func ExampleTail_tail1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	result, err := Tail(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "grayson", "tim" }
}

func ExampleTail_tail2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result, err := Tail(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 2, 3, 4, 5, 6, 7, 8, 9 }
}

func ExampleTake_take1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	take := 2

	result, err := Take(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "jason", "damian" }
}

func ExampleTake_take2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	take := 5

	result, err := Take(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 5 }
}

func ExampleTakeRight_takeRight1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	take := 2

	result, err := TakeRight(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "grayson", "tim" }
}

func ExampleTakeRight_takeRight2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	take := 5

	result, err := TakeRight(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 5, 6, 7, 8, 9 }
}

func ExampleUnion_union1() {
	result, err := Union(
		[]string{"damian", "grayson", "grayson", "cassandra"},
		[]string{"tim", "grayson", "jason", "stephanie"},
		[]string{"duke"},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "grayson", "cassandra", "tim", "jason", "stephanie", "duke" }
}

func ExampleUnion_union2() {
	result, err := Union(
		[]int{1, 2, 3},
		[]int{2, 3, 4, 5, 6},
		[]int{2, 5, 7, 8},
		[]int{9},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 5, 6, 7, 8, 9 }
}

func ExampleUniq_uniq1() {
	data := []string{"damian", "grayson", "grayson", "cassandra"}
	result, err := Uniq(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "grayson", "cassandra" }
}

func ExampleUniq_uniq2() {
	data := []float64{1.1, 3.00000, 3.1, 2.2000000, 3, 2.2, 3.0}
	result, err := Uniq(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 3, 3.1, 2.2 }
}

func ExampleWithout_without1() {
	data := []int{1, 2, 3, 4, 5, 6}
	exclude := []int{3}

	result, err := Without(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 4, 5, 6 }
}

func ExampleWithout_without2() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}
	exclude := []float64{2.1, 3.2, 6.3}

	result, err := Without(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 4.2, 5.2 }
}

func ExampleWithout_without3() {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	exclude := []string{"grayson", "tim"}

	result, err := Without(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "cassandra", "jason", "stephanie" }
}
