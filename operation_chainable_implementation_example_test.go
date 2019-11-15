package gubrak

import (
	"fmt"
	"log"
	"strings"
)

func ExampleChainable_Chunk_chunk1() {
	data := []int{1, 2, 3, 4, 5}
	size := 2

	result := From(data).Chunk(size).Result()
	fmt.Println(result)
	// ===> [][]int{ { 1, 2 }, { 3, 4 }, { 5 } }
}

func ExampleChainable_Chunk_chunk2() {
	data := []string{"a", "b", "c", "d", "e"}
	size := 3

	result := From(data).Chunk(size).Result()
	fmt.Println(result)
	// ===> [][]string{ { "a", "b", "c" }, { "d", "e" } }
}

func ExampleChainable_Chunk_chunk3() {
	data := []interface{}{
		3.2, "a", -1,
		make([]byte, 0),
		map[string]int{"b": 2},
		[]string{"a", "b", "c"},
	}
	size := 3

	result, err := From(data).Chunk(size).ResultAndError()
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

func ExampleChainable_Compact_compact1() {
	data := []int{-2, -1, 0, 1, 2}

	result := From(data).Compact().Result()
	fmt.Println(result)
	// ===> []int{ -2, -1, 1, 2 }
}

func ExampleChainable_Compact_compact2() {
	data := []string{"a", "b", "", "d"}

	result := From(data).Compact().Result()
	fmt.Println(result)
	// ===> []string{ "a", "b", "d" }
}

func ExampleChainable_Compact_compact3() {
	data := []interface{}{-2, 0, 1, 2, false, true, "", "hello", nil}

	result := From(data).Compact().Result()
	fmt.Println(result)
	// ===> []interface{}{ -2, 1, 2, true, "hello" }
}

func ExampleChainable_Compact_compact4() {
	item1, item2, item3 := "a", "b", "c"
	data := []*string{&item1, nil, &item2, nil, &item3}

	result, err := From(data).Compact().ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []*string{ (*string)(0xc42000e1e0), (*string)(0xc42000e1f0), (*string)(0xc42000e200) }
}

func ExampleChainable_ConcatMany_concat1() {
	data := []int{1, 2, 3, 4}
	dataConcat1 := []int{4, 6, 7}
	dataConcat2 := []int{8, 9}

	result := From(data).ConcatMany(dataConcat1, dataConcat2).Result()
	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 5, 6, 7, 8, 9 }
}

func ExampleChainable_ConcatMany_concat2() {
	data := []string{"my"}
	dataConcat1 := []string{"name", "is"}
	dataConcat2 := []string{"jason", "todd"}

	result, err := From(data).ConcatMany(dataConcat1, dataConcat2).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "my", "name", "is", "jason", "todd" }
}

func ExampleChainable_Concat_concat1() {
	data := []int{1, 2, 3, 4}
	dataConcat := []int{4, 6, 7}

	result := From(data).Concat(dataConcat).Result()
	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 5, 6, 7 }
}

func ExampleChainable_Count_countMap1() {
	data := map[string]interface{}{
		"name":   "jason",
		"age":    12,
		"isMale": true,
	}
	result := From(data).Count().Result()
	fmt.Println(result)
	// ===> 3
}

func ExampleChainable_CountBy_countMap2() {
	data := map[string]interface{}{
		"name":   "jason",
		"age":    12,
		"isMale": true,
	}

	result := From(data).
		CountBy(func(val interface{}, key string) bool {
			return strings.Contains(strings.ToLower(key), "m")
		}).
		Result()
	fmt.Println(result)
	// ===> 2
}

func ExampleChainable_CountBy_countMap3() {
	data := map[string]interface{}{
		"name":   "jason",
		"age":    12,
		"isMale": true,
	}

	result := From(data).
		CountBy(func(val interface{}, key string, i int) bool {
			return strings.Contains(strings.ToLower(key), "m") && i > 1
		}).
		Result()
	fmt.Println(result)
	// ===> 1
}

func ExampleChainable_Count_countSlice1() {
	data := []string{"damian", "grayson", "cassandra"}
	result := From(data).Count().Result()
	fmt.Println(result)
	// ===> 3
}

func ExampleChainable_CountBy_countSlice2() {
	data := []string{"damian", "grayson", "cassandra"}

	result := From(data).
		CountBy(func(each string) bool {
			return strings.Contains(each, "d")
		}).
		Result()
	fmt.Println(result)
	// ===> 2
}

func ExampleChainable_CountBy_countSlice3() {
	data := []string{"damian", "grayson", "cassandra"}

	result, err := From(data).
		CountBy(func(each string, i int) bool {
			return len(each) > 6 && i > 1
		}).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 1
}

func ExampleChainable_Difference_difference1() {
	data := []int{1, 2, 3, 4, 4, 6, 7}
	dataDiff := []int{2, 7}

	result, err := From(data).Difference(dataDiff).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 3, 4, 4, 6 }
}

func ExampleChainable_DifferenceMany_difference2() {
	data := []string{"a", "b", "b", "c", "d", "e", "f", "g", "h"}
	dataDiff1 := []string{"b", "d"}
	dataDiff2 := []string{"e", "f", "h"}

	result, err := From(data).DifferenceMany(dataDiff1, dataDiff2).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "a", "c", "g" }
}

func ExampleChainable_DifferenceMany_difference3() {
	data := []float64{1.1, 1.11, 1.2, 2.3, 3.0, 3, 4.0, 4.00000, 4.000000001}
	dataDiff1 := []float64{1.1, 3}
	dataDiff2 := []float64{4.000000001}

	result, err := From(data).DifferenceMany(dataDiff1, dataDiff2).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.11, 1.2, 2.3, 4, 4 }
}

func ExampleChainable_Drop_drop1() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	n := 1

	result, err := From(data).Drop(n).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 2, 3, 4, 4, 5, 6 }
}

func ExampleChainable_Drop_drop2() {
	data := []string{"a", "b", "c", "d", "e", "f"}
	n := 3

	result, err := From(data).Drop(n).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "d", "e", "f" }
}

func ExampleChainable_DropRight_dropRight1() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	n := 1

	result, err := From(data).Drop(n).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 4, 5 }
}

func ExampleChainable_DropRight_dropRight2() {
	data := []string{"a", "b", "c", "d", "e", "f"}
	n := 3

	result, err := From(data).Drop(n).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "a", "b", "c" }
}

func ExampleChainable_Each_eachMap1() {
	data := map[string]interface{}{
		"name":   "damian",
		"age":    17,
		"gender": "male",
	}

	err := From(data).
		Each(func(value interface{}, key string) {
			fmt.Printf("%s: %v \n", key, value)
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_Each_eachMap2() {
	data := map[string]interface{}{
		"name":   "damian",
		"age":    17,
		"gender": "male",
	}

	err := From(data).
		Each(func(value interface{}, key string, i int) {
			fmt.Printf("key: %s, value: %v, index: %d \n", key, value, i)
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_Each_eachSlice1() {
	data := []string{"damian", "grayson", "cassandra"}

	err := From(data).
		Each(func(each string) {
			fmt.Println(each)
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_Each_eachSlice2() {
	data := []string{"damian", "grayson", "cassandra"}

	err := From(data).
		Each(func(each string, i int) {
			fmt.Printf("element %d: %s \n", i, each)
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_Each_eachSlice3() {
	type Sample struct {
		Name string
		Age  int
	}

	data := []Sample{
		{Name: "damian", Age: 12},
		{Name: "grayson", Age: 10},
		{Name: "cassandra", Age: 11},
	}

	err := From(data).
		Each(func(each Sample) {
			fmt.Printf("name: %s, age: %d \n", each.Name, each.Age)
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_Each_eachSlice4() {
	data := []string{"damian", "grayson", "cassandra", "tim", "jason", "stephanie"}

	err := From(data).
		Each(func(each string, i int) bool {
			if i > 3 { // will stop after fourth loop
				return false
			}

			fmt.Println(each)
			return true
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_EachRight_eachRightMap1() {
	data := map[string]interface{}{
		"name":   "damian",
		"age":    17,
		"gender": "male",
	}

	err := From(data).
		EachRight(func(value interface{}, key string) {
			fmt.Printf("%s: %v \n", key, value)
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_EachRight_eachRightMap2() {
	data := map[string]interface{}{
		"name":   "damian",
		"age":    17,
		"gender": "male",
	}

	err := From(data).
		EachRight(func(value interface{}, key string, i int) {
			fmt.Printf("key: %s, value: %v, index: %d \n", key, value, i)
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_EachRight_eachRightSlice1() {
	data := []string{"damian", "grayson", "cassandra"}

	err := From(data).
		EachRight(func(each string) {
			fmt.Println(each)
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_EachRight_eachRightSlice2() {
	data := []string{"damian", "grayson", "cassandra"}

	err := From(data).
		EachRight(func(each string, i int) {
			fmt.Printf("element %d: %s \n", i, each)
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_EachRight_eachRightSlice3() {
	type Sample struct {
		Name string
		Age  int
	}

	data := []Sample{
		{Name: "damian", Age: 12},
		{Name: "grayson", Age: 10},
		{Name: "cassandra", Age: 11},
	}

	err := From(data).
		EachRight(func(each Sample) {
			fmt.Printf("name: %s, age: %d \n", each.Name, each.Age)
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_EachRight_eachRightSlice4() {
	data := []string{"damian", "grayson", "cassandra", "tim", "jason", "stephanie"}

	err := From(data).
		EachRight(func(each string, i int) bool {
			if i > 3 { // will stop after fourth loop
				return false
			}

			fmt.Println(each)
			return true
		}).
		Error()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleChainable_Fill_fill1() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	replacement := 9

	result, err := From(data).Fill(replacement).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 9, 9, 9, 9, 9, 9, 9 }
}

func ExampleChainable_Fill_fill2() {
	data := []string{"grayson", "jason", "tim", "damian"}
	replacement := "alfred"
	start := 2

	result, err := From(data).Fill(replacement, start).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ "grayson", "jason", "alfred", "alfred" }
}

func ExampleChainable_Fill_fill3() {
	data := []float64{1, 2.2, 3.0002, 4, 4, 5.12, 6}
	replacement := float64(9)
	start, end := 3, 5

	result, err := From(data).Fill(replacement, start, end).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1, 2.2, 3.0002, 9, 9, 5.12, 6 }
}

func ExampleChainable_Filter_filterMap() {
	data := map[string]int{
		"clean code":       10000,
		"rework":           12000,
		"detective comics": 11500,
	}

	result, err := From(data).
		Filter(func(value int, key string) bool {
			return value > 11000
		}).
		ResultAndError()
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

func ExampleChainable_Filter_filterSlice() {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := From(data).
		Filter(func(each Sample) bool {
			return each.DailyDownloads > 11000
		}).
		ResultAndError()
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

func ExampleChainable_Find_find1() {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := From(data).
		Find(func(each Sample) bool {
			return each.EbookName == "rework"
		}).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> Sample { EbookName: "rework", DailyDownloads: 12000 }
}

func ExampleChainable_Find_find2() {
	data := []string{"clean code", "rework", "detective comics"}

	result, err := From(data).
		Find(func(each string, i int) bool {
			return strings.Contains(each, "co")
		}, 1).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "detective comics"
}

func ExampleChainable_FindIndex_findIndex1() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}
	predicate := func(each string) bool {
		return each == "tim"
	}

	result, err := From(data).
		FindIndex(predicate).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 3
}

func ExampleChainable_FindIndex_findIndex2() {
	data := []int{-2, -1, 0, 1, 2}

	result, err := From(data).
		FindIndex(func(each int) bool {
			return each == 4
		}).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> -1
}

func ExampleChainable_FindIndex_findIndex3() {
	data := []float64{1, 1.1, 1.2, 1.200001, 1.2000000001, 1.3}

	result, err := From(data).
		FindIndex(func(each float64) bool {
			return each == 1.2000000001
		}).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 4
}

func ExampleChainable_FindIndex_findIndex4() {
	data := []int{1, 2, 3, 3, 4, 5}
	predicate := func(each int) bool {
		return each == 3
	}
	fromIndex := 2

	result, err := From(data).
		FindIndex(predicate, fromIndex).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 2
}

func ExampleChainable_FindIndex_findIndex5() {
	data := []int{1, 2, 3, 3, 4, 5}
	predicate := func(each int) bool {
		return each == 3
	}
	fromIndex := 3

	result, err := From(data).
		FindIndex(predicate, fromIndex).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 3
}

func ExampleChainable_FindLast_findLast1() {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := From(data).
		FindLast(func(each Sample) bool {
			return strings.Contains(each.EbookName, "co")
		}).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> Sample { EbookName: "detective comics", DailyDownloads: 11500 }
}

func ExampleChainable_FindLast_findLast2() {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	result, err := From(data).
		FindLast(func(each string, i int) bool {
			return strings.Contains(each, "co")
		}, 2).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "detective comics"
}

func ExampleChainable_FindLast_findLast3() {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	result, err := From(data).
		FindLast(func(each string, i int) bool {
			return strings.Contains(each, "co")
		}, 3).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "coco"
}

func ExampleChainable_FindLastIndex_findLastIndex1() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}

	result, err := From(data).
		FindLastIndex(func(each string) bool {
			return each == "tim"
		}).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 4
}

func ExampleChainable_FindLastIndex_findLastIndex2() {
	data := []int{1, 2, 2, 3, 3, 4, 5}
	predicate := func(each int) bool {
		return each == 3
	}
	fromIndex := 4

	result, err := From(data).
		FindLastIndex(predicate, fromIndex).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 4
}

func ExampleChainable_FindLastIndex_findLastIndex3() {
	data := []int{1, 2, 3, 3, 4, 5}
	predicate := func(each int) bool {
		return each == 3
	}
	fromIndex := 3

	result, err := From(data).
		FindLastIndex(predicate, fromIndex).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 3
}

func ExampleChainable_FindLastIndex_findLastIndex4() {
	data := []int{1, 2, 3, 3, 4, 5}
	predicate := func(each int) bool {
		return each == 3
	}
	fromIndex := 2

	result, err := From(data).
		FindLastIndex(predicate, fromIndex).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> -1
}

func ExampleChainable_First_first1() {
	data := []string{"damian", "grayson", "cassandra"}
	result := From(data).First().Result()

	fmt.Println(result)
	// ===> "damian"
}

func ExampleChainable_First_first2() {
	data := []string{}
	result := From(data).First().Result()

	fmt.Println(result)
	// ===> nil
}

func ExampleChainable_FromPairs_fromPairs1() {
	data := []interface{}{
		[]interface{}{"a", 1},
		[]interface{}{"b", 2},
	}

	result, err := From(data).FromPairs().ResultAndError()
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

func ExampleChainable_FromPairs_fromPairs2() {
	data := []interface{}{
		[]interface{}{true, []int{1, 2, 3}},
		[]interface{}{false, []string{"damian", "grayson"}},
	}

	result, err := From(data).FromPairs().ResultAndError()
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

func ExampleChainable_GroupBy_groupBy1() {
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

	result, err := From(data).
		GroupBy(func(each Sample) string {
			return each.Category
		}).
		ResultAndError()
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

func ExampleChainable_GroupBy_groupBy2() {
	data := []int{1, 2, 3, 5, 6, 4, 2, 5, 2}

	result, err := From(data).
		GroupBy(func(each int) int {
			return each
		}).
		ResultAndError()
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

func ExampleChainable_Contains_includesMap1() {
	data := map[string]string{
		"name":  "grayson",
		"hobby": "helping people",
	}

	result := From(data).Contains("grayson").Result()

	fmt.Println(result)
	// ===> true
}

func ExampleChainable_Contains_includesMap2() {
	data := map[string]string{
		"name":  "grayson",
		"hobby": "helping people",
	}

	result := From(data).Contains("batmobile").Result()

	fmt.Println(result)
	// ===> false
}

func ExampleChainable_Contains_includesSlice1() {
	data := []string{"damian", "tim", "jason", "grayson"}

	result := From(data).Contains("tim").Result()

	fmt.Println(result)
	// ===> true
}

func ExampleChainable_Contains_includesSlice2() {
	data := []string{"damian", "tim", "jason", "grayson"}

	result := From(data).Contains("tim", 2).Result()

	fmt.Println(result)
	// ===> false
}

func ExampleChainable_Contains_includesSlice3() {
	data := []string{"damian", "tim", "jason", "grayson"}

	result := From(data).Contains("cassandra").Result()

	fmt.Println(result)
	// ===> false
}

func ExampleChainable_Contains_includesSlice4() {
	data := []interface{}{"name", 12, true}

	From(data).Contains("name").Result() // ===> true
	From(data).Contains(12).Result()     // ===> true
	From(data).Contains(true).Result()   // ===> true
}

func ExampleChainable_Contains_includesSlice5() {
	From("damian").Contains("an").Result() // ===> true
}

func ExampleChainable_IndexOf_indexOf1() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}

	From(data).IndexOf("duke").Result()    // ===> -1
	From(data).IndexOf("tim").Result()     // ===> 3
	From(data).IndexOf("tim", 4).Result()  // ===> 4
	From(data).IndexOf("tim", -4).Result() // ===> 3
	From(data).IndexOf("tim", -3).Result() // ===> 4
	From(data).IndexOf("tim", -2).Result() // ===> -1
}

func ExampleChainable_IndexOf_indexOf2() {
	data := []float64{2.1, 2.2, 3, 3.00000, 3.1, 3.9, 3.95}

	From(data).IndexOf(2.2).Result()           // ===> 1
	From(data).IndexOf(3).Result()             // ===> -1
	From(data).IndexOf(float64(3)).Result()    // ===> 2 (because 3 is detected as int32, not float64)
	From(data).IndexOf(float64(3), 2).Result() // ===> 2
	From(data).IndexOf(float64(3), 3).Result() // ===> 3
}

func ExampleChainable_IndexOf_indexOf3() {
	data := []interface{}{"jason", 24, true}

	From(data).IndexOf(24).Result()     // ===> 1
	From(data).IndexOf(24, -1).Result() // ===> -1
}

func ExampleChainable_Initial_initial1() {
	data := []string{"damian", "grayson", "cassandra"}

	result := From(data).Initial().Result()

	fmt.Println(result)
	// ===> []string{ "damian", "grayson" }
}

func ExampleChainable_Initial_initial2() {
	data := []int{1, 2, 3, 4, 5}

	result := From(data).Initial().Result()

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4 }
}

func ExampleChainable_Initial_initial3() {
	data := []map[string]string{{"name": "jason"}}

	result := From(data).Initial().Result()

	fmt.Println(result)
	// ===> []map[string]string{}
}

func ExampleChainable_Initial_initial4() {
	data := []float64{}

	result := From(data).Initial().Result()

	fmt.Println(result)
	// ===> []float64{}
}

func ExampleChainable_Intersection_intersection1() {
	result := From([]string{"damian", "grayson", "cassandra", "tim", "tim", "jason"}).
		Intersection([]string{"cassandra", "tim", "jason"}).
		Result()

	fmt.Println(result)
	// ===> []string{ "cassandra", "tim", "jason" }
}

func ExampleChainable_Intersection_intersection2() {
	result := From([]float64{0.8, 0.8001, 0.999, 1, 1.0, 1.000001, 1.1000000, 1.1001, 1.2, 1.33, 1.4}).
		IntersectionMany(
			[]float64{0.8, 0.8001, 0.999, 1, 1.0, 1.000001, 1.1000000, 1.2, 1.33},
			[]float64{1.1000000, 1.2, 0.8001, 0.999, 1.33, 1, 1.0, 1.000001},
			[]float64{1.2, 0.8001, 0.999, 1.33, 1.000092},
			[]float64{0.8001, 0.999, 1.33, 1.400001},
		).
		Result()

	fmt.Println(result)
	// ===> []float64{ 0.8001, 0.999, 1.33 }
}

func ExampleChainable_Join_join1() {
	data := []string{"damian", "grayson", "cassandra"}
	separator := " - "

	result := From(data).Join(separator).Result()

	fmt.Println(result)
	// ===> "damian - grayson - cassandra"
}

func ExampleChainable_Join_join2() {
	data := []int{1, 2, 3, 4}
	separator := ", "

	result := From(data).Join(separator).Result()

	fmt.Println(result)
	// ===> "1, 2, 3, 4"
}

func ExampleChainable_KeyBy_keyBy() {
	type HashMap map[string]string

	data := []HashMap{
		{"name": "grayson", "hobby": "helping people"},
		{"name": "jason", "hobby": "punching people"},
		{"name": "tim", "hobby": "stay awake all the time"},
		{"name": "damian", "hobby": "getting angry"},
	}

	result, err := From(data).
		KeyBy(func(each HashMap) string {
			return each["name"]
		}).
		ResultAndError()

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

func ExampleChainable_Last_last1() {
	data := []string{"damian", "grayson", "cassandra"}

	result := From(data).Last().Result()

	fmt.Println(result)
	// ===> "cassandra"
}

func ExampleChainable_Last_last2() {
	data := []int{1}

	result := From(data).Last().Result()

	fmt.Println(result)
	// ===> 1
}

func ExampleChainable_Last_last3() {
	data := []string{}

	result := From(data).Last().Result()

	fmt.Println(result)
	// ===> nil
}

func ExampleChainable_LastIndexOf_lastIndexOf1() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}

	From(data).LastIndexOf("duke").Result()    // ===> -1
	From(data).LastIndexOf("tim").Result()     // ===> 4
	From(data).LastIndexOf("tim", 4).Result()  // ===> 4
	From(data).LastIndexOf("tim", -4).Result() // ===> 3
	From(data).LastIndexOf("tim", -3).Result() // ===> 4
	From(data).LastIndexOf("tim", -2).Result() // ===> 4
}

func ExampleChainable_LastIndexOf_lastIndexOf2() {
	data := []float64{2.1, 2.2, 3, 3.00000, 3.1, 3.9, 3.95}

	From(data).LastIndexOf(2.2).Result()           // ===> 1
	From(data).LastIndexOf(3).Result()             // ===> -1 (because 3 is detected as int32, not float64)
	From(data).LastIndexOf(float64(3)).Result()    // ===> 3
	From(data).LastIndexOf(float64(3), 2).Result() // ===> 2
	From(data).LastIndexOf(float64(3), 3).Result() // ===> 3
}

func ExampleChainable_LastIndexOf_lastIndexOf3() {
	data := []interface{}{"jason", 24, true}

	From(data).LastIndexOf(24).Result()     // ===> 1
	From(data).LastIndexOf(24, -1).Result() // ===> 1
}

func ExampleChainable_Map_map1() {
	type Sample struct {
		EbookName      string
		DailyDownloads int
	}

	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result, err := From(data).
		Map(func(each Sample, i int) string {
			return each.EbookName
		}).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "clean code", "rework", "detective comics" }
}

func ExampleChainable_Map_map2() {
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

	result, err := From(data).
		Map(func(each SampleOne, i int) SampleTwo {
			ebook := each.EbookName
			if !each.IsActive {
				ebook = fmt.Sprintf("%s (inactive)", each.EbookName)
			}

			downloadsInThousands := float32(each.DailyDownloads) / float32(1000)

			return SampleTwo{Ebook: ebook, DownloadsInThousands: downloadsInThousands}
		}).
		ResultAndError()
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

func ExampleChainable_Nth_nth1() {
	data := []string{"grayson", "jason", "tim", "damian"}

	From(data).Nth(1).Result()  // ===> "jason"
	From(data).Nth(2).Result()  // ===> "tim"
	From(data).Nth(-1).Result() // ===> "damian"
}

func ExampleChainable_Nth_nth2() {
	data := []int{1, 2, 3, 4, 5}
	From(data).Nth(4).Result() // ===> 5
}

func ExampleChainable_OrderBy_orderBy1() {
	type HashMap map[string]string

	data := []HashMap{
		{"name": "tim", "hobby": "stay awake all the time"},
		{"name": "grayson", "hobby": "helping people"},
		{"name": "damian", "hobby": "getting angry"},
		{"name": "jason", "hobby": "punching people"},
	}

	result, err := From(data).
		OrderBy(func(each HashMap) string {
			return each["name"]
		}).
		ResultAndError()
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

func ExampleChainable_OrderBy_orderBy2() {
	type HashMap map[string]interface{}

	data := []HashMap{
		{"name": "tim", "hobby": "stay awake all the time", "age": 20},
		{"name": "grayson", "hobby": "helping people", "age": 24},
		{"name": "damian", "hobby": "getting angry", "age": 17},
		{"name": "jason", "hobby": "punching people", "age": 22},
	}

	result, err := From(data).
		OrderBy(func(each HashMap) int {
			return each["age"].(int)
		}).
		ResultAndError()
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

func ExampleChainable_OrderBy_orderBy3() {
	type HashMap map[string]interface{}

	data := []HashMap{
		{"name": "tim", "hobby": "stay awake all the time", "age": 20},
		{"name": "grayson", "hobby": "helping people", "age": 24},
		{"name": "damian", "hobby": "getting angry", "age": 17},
		{"name": "jason", "hobby": "punching people", "age": 22},
	}

	result, err := From(data).
		OrderBy(func(each HashMap) int {
			return each["age"].(int)
		}, false).
		ResultAndError()
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

func ExampleChainable_OrderBy_orderBy4() {
	type HashMap map[string]interface{}

	data := []HashMap{
		{"name": "tim", "hobby": "stay awake all the time", "age": 20},
		{"name": "grayson", "hobby": "helping people", "age": 24},
		{"name": "damian", "hobby": "getting angry", "age": 17},
		{"name": "jason", "hobby": "punching people", "age": 22},
	}

	result, err := From(data).
		OrderBy(func(each HashMap) int {
			return each["age"].(int)
		}, true, false).
		ResultAndError()
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

func ExampleChainable_Partition_partition() {
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

	resultTruthy, resultFalsey, err := From(data).
		Partition(func(each HashMap) bool {
			return each["isMale"].(bool)
		}).
		ResultAndError()
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

func ExampleChainable_Exclude_exclude1() {
	data := []int{1, 2, 3, 4, 5, 6}
	result, err := From(data).Exclude(3).ResultAndError()

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 4, 5, 6 }
}

func ExampleChainable_ExcludeMany_excludeMany1() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}
	result, err := From(data).ExcludeMany(2.1, 3.2, 6.3).ResultAndError()

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 4.2, 5.2 }
}

func ExampleChainable_ExcludeMany_excludeMany2() {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := From(data).ExcludeMany("grayson", "tim").ResultAndError()

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "cassandra", "jason", "stephanie" }
}

func ExampleChainable_ExcludeAt_excludeAt() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}

	result, err := From(data).ExcludeAt(1).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 2.1, 3.2, 4.2, 5.2, 6.3 }
}

func ExampleChainable_ExcludeAt_excludeAtMany() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}

	result, err := From(data).ExcludeAtMany(1, 3).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 3.2, 5.2, 6.3 }
}

func ExampleChainable_Reduce_reduceMap1() {
	type HashMap map[string]interface{}

	data := HashMap{
		"name":   "grayson",
		"age":    21,
		"isMale": true,
	}

	result, err := From(data).
		Reduce(func(accumulator string, value interface{}, key string) string {
			if accumulator == "" {
				accumulator = fmt.Sprintf("%s: %v", key, value)
			} else {
				accumulator = fmt.Sprintf("%s, %s: %v", accumulator, key, value)
			}

			return accumulator
		}, "").
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "name: grayson, age: 21, isMale: true"
}

func ExampleChainable_Reduce_reduceSlice1() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result, err := From(data).
		Reduce(func(accumulator, each int) int {
			return accumulator + each
		}, 0).
		ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 55
}

func ExampleChainable_Reduce_reduceSlice2() {
	type HashMap map[string]interface{}

	data := [][]interface{}{
		{"name", "grayson"},
		{"age", 21},
		{"isMale", true},
	}

	result, err := From(data).
		Reduce(func(accumulator HashMap, each []interface{}, i int) HashMap {
			accumulator[each[0].(string)] = each[1]
			return accumulator
		}, HashMap{}).
		ResultAndError()
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

func ExampleChainable_Reject_rejectMap() {
	data := map[string]int{
		"clean code":       10000,
		"rework":           12000,
		"detective comics": 11500,
	}

	result, err := From(data).Reject(func(value int, key string) bool {
		return value > 11000
	}).ResultAndError()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> map[string]int{ "clean code": 10000 }
}

func ExampleChainable_Reject_rejectSlice() {
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

func ExampleChainable_Remove_remove1() {
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

func ExampleChainable_Remove_remove2() {
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

func ExampleChainable_Reverse_reverse1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	result, err := Reverse(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "tim", "grayson", "damian", "jason" }
}

func ExampleChainable_Reverse_reverse2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result, err := Reverse(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 9, 8, 7, 6, 5, 4, 3, 2, 1 }
}

func ExampleChainable_Sample_sample() {
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

func ExampleChainable_SampleSize_sampleSize() {
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

func ExampleChainable_Shuffle_shuffle1() {
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

func ExampleChainable_Shuffle_shuffle2() {
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

func ExampleChainable_Size_sizeSlice() {
	Size([]int{1, 2, 3, 4, 5}) // ===> 5
	Size("bruce")              // ===> 5
}

func ExampleChainable_Size_sizeMap() {
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

func ExampleChainable_Tail_tail1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	result, err := Tail(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "grayson", "tim" }
}

func ExampleChainable_Tail_tail2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result, err := Tail(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 2, 3, 4, 5, 6, 7, 8, 9 }
}

func ExampleChainable_Take_take1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	take := 2

	result, err := Take(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "jason", "damian" }
}

func ExampleChainable_Take_take2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	take := 5

	result, err := Take(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 5 }
}

func ExampleChainable_TakeRight_takeRight1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	take := 2

	result, err := TakeRight(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "grayson", "tim" }
}

func ExampleChainable_TakeRight_takeRight2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	take := 5

	result, err := TakeRight(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 5, 6, 7, 8, 9 }
}

func ExampleChainable_Union_union1() {
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

func ExampleChainable_Union_union2() {
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

func ExampleChainable_Uniq_uniq1() {
	data := []string{"damian", "grayson", "grayson", "cassandra"}
	result, err := Uniq(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "grayson", "cassandra" }
}

func ExampleChainable_Uniq_uniq2() {
	data := []float64{1.1, 3.00000, 3.1, 2.2000000, 3, 2.2, 3.0}
	result, err := Uniq(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 3, 3.1, 2.2 }
}

func ExampleChainable_Without_without1() {
	data := []int{1, 2, 3, 4, 5, 6}
	exclude := []int{3}

	result, err := Without(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 4, 5, 6 }
}

func ExampleChainable_Without_without2() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}
	exclude := []float64{2.1, 3.2, 6.3}

	result, err := Without(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 4.2, 5.2 }
}

func ExampleChainable_Without_without3() {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	exclude := []string{"grayson", "tim"}

	result, err := Without(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "cassandra", "jason", "stephanie" }
}
