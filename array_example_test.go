package gubrak

import (
	"fmt"
	"log"
	"strings"
)

func ExampleChunk_1() {
	data := []int{1, 2, 3, 4, 5}
	size := 2

	result, err := Chunk(data, size)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> [][]int{ { 1, 2 }, { 3, 4 }, { 5 } }
}

func ExampleChunk_2() {
	data := []string{"a", "b", "c", "d", "e"}
	size := 3

	result, err := Chunk(data, size)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> [][]string{ { "a", "b", "c" }, { "d", "e" } }
}

func ExampleChunk_3() {
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

func ExampleCompact_1() {
	data := []int{-2, -1, 0, 1, 2}

	result, err := Compact(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ -2, -1, 1, 2 }
}

func ExampleCompact_2() {
	data := []string{"a", "b", "", "d"}

	result, err := Compact(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "a", "b", "d" }
}

func ExampleCompact_3() {
	data := []interface{}{-2, 0, 1, 2, false, true, "", "hello", nil}

	result, err := Compact(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []interface{}{ -2, 1, 2, true, "hello" }
}

func ExampleCompact_4() {
	item1, item2, item3 := "a", "b", "c"
	data := []*string{&item1, nil, &item2, nil, &item3}

	result, err := Compact(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []*string{ (*string)(0xc42000e1e0), (*string)(0xc42000e1f0), (*string)(0xc42000e200) }
}

func ExampleConcat_1() {
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

func ExampleConcat_2() {
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

func ExampleDifference_1() {
	data := []int{1, 2, 3, 4, 4, 6, 7}
	dataDiff := []int{2, 7}

	result, err := Difference(data, dataDiff)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 3, 4, 4, 6 }
}

func ExampleDifference_2() {
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

func ExampleDifference_3() {
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

func ExampleDrop_1() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	n := 1

	result, err := Drop(data, n)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 2, 3, 4, 4, 5, 6 }
}

func ExampleDrop_2() {
	data := []string{"a", "b", "c", "d", "e", "f"}
	n := 3

	result, err := Drop(data, n)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "d", "e", "f" }
}

func ExampleDropRight_1() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	n := 1

	result, err := DropRight(data, n)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 4, 5 }
}

func ExampleDropRight_2() {
	data := []string{"a", "b", "c", "d", "e", "f"}
	n := 3

	result, err := DropRight(data, n)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "a", "b", "c" }
}

func ExampleFill_1() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	replacement := 9

	result, err := Fill(data, replacement)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 9, 9, 9, 9, 9, 9, 9 }
}

func ExampleFill_2() {
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

func ExampleFill_3() {
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

func ExampleFindIndex_1() {
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

func ExampleFindIndex_2() {
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

func ExampleFindIndex_3() {
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

func ExampleFindIndex_4() {
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

func ExampleFindIndex_5() {
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

func ExampleFindLastIndex_1() {
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

func ExampleFindLastIndex_2() {
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

func ExampleFindLastIndex_3() {
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

func ExampleFindLastIndex_4() {
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

func ExampleFirst_1() {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := First(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "damian"
}

func ExampleFirst_2() {
	data := []string{}
	result, err := First(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> nil
}

func ExampleFromPairs_1() {
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

func ExampleFromPairs_2() {
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

func ExampleIndexOf_1() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}
	IndexOf(data, "duke")    // ===> -1
	IndexOf(data, "tim")     // ===> 3
	IndexOf(data, "tim", 4)  // ===> 4
	IndexOf(data, "tim", -4) // ===> 3
	IndexOf(data, "tim", -3) // ===> 4
	IndexOf(data, "tim", -2) // ===> -1
}

func ExampleIndexOf_2() {
	data := []float64{2.1, 2.2, 3, 3.00000, 3.1, 3.9, 3.95}
	IndexOf(data, 2.2)           // ===> 1
	IndexOf(data, 3)             // ===> -1
	IndexOf(data, float64(3))    // ===> 2 (because 3 is detected as int32, not float64)
	IndexOf(data, float64(3), 2) // ===> 2
	IndexOf(data, float64(3), 3) // ===> 3
}

func ExampleIndexOf_3() {
	data := []interface{}{"jason", 24, true}
	IndexOf(data, 24)     // ===> 1
	IndexOf(data, 24, -1) // ===> -1
}

func ExampleInitial_1() {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Initial(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "grayson" }
}

func ExampleInitial_2() {
	data := []int{1, 2, 3, 4, 5}
	result, err := Initial(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4 }
}

func ExampleInitial_3() {
	data := []map[string]string{{"name": "jason"}}
	result, err := Initial(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []map[string]string{}
}

func ExampleInitial_4() {
	data := []float64{}
	result, err := Initial(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{}
}

func ExampleIntersection_1() {
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

func ExampleIntersection_2() {
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

func ExampleJoin_1() {
	data := []string{"damian", "grayson", "cassandra"}
	separator := " - "

	result, err := Join(data, separator)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "damian - grayson - cassandra"
}

func ExampleJoin_2() {
	data := []int{1, 2, 3, 4}
	separator := ", "

	result, err := Join(data, separator)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "1, 2, 3, 4"
}

func ExampleLast_1() {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Last(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "cassandra"
}

func ExampleLast_2() {
	data := []int{1}
	result, err := Last(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 1
}

func ExampleLast_3() {
	data := []string{}
	result, err := Last(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> nil
}

func ExampleLastIndexOf_1() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}

	LastIndexOf(data, "duke")    // ===> -1
	LastIndexOf(data, "tim")     // ===> 4
	LastIndexOf(data, "tim", 4)  // ===> 4
	LastIndexOf(data, "tim", -4) // ===> 3
	LastIndexOf(data, "tim", -3) // ===> 4
	LastIndexOf(data, "tim", -2) // ===> 4
}

func ExampleLastIndexOf_2() {
	data := []float64{2.1, 2.2, 3, 3.00000, 3.1, 3.9, 3.95}

	LastIndexOf(data, 2.2)           // ===> 1
	LastIndexOf(data, 3)             // ===> -1 (because 3 is detected as int32, not float64)
	LastIndexOf(data, float64(3))    // ===> 3
	LastIndexOf(data, float64(3), 2) // ===> 2
	LastIndexOf(data, float64(3), 3) // ===> 3
}

func ExampleLastIndexOf_3() {
	data := []interface{}{"jason", 24, true}

	LastIndexOf(data, 24)     // ===> 1
	LastIndexOf(data, 24, -1) // ===> 1
}

func ExampleNth_1() {
	data := []string{"grayson", "jason", "tim", "damian"}

	Nth(data, 1)  // ===> "jason"
	Nth(data, 2)  // ===> "tim"
	Nth(data, -1) // ===> "damian"
}

func ExampleNth_2() {
	data := []int{1, 2, 3, 4, 5}
	result, err := Nth(data, 4)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 5
}

func ExamplePull_1() {
	data := []int{1, 2, 3, 4, 5, 6}
	result, err := Pull(data, 3)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 4, 5, 6 }
}

func ExamplePull_2() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}
	result, err := Pull(data, 2.1, 3.2, 6.3)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 4.2, 5.2 }
}

func ExamplePull_3() {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	result, err := Pull(data, "grayson", "tim")

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "cassandra", "jason", "stephanie" }
}

func ExamplePullAll_1() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}
	exclude := []float64{2.1, 3.2, 6.3}

	result, err := PullAll(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 4.2, 5.2 }
}

func ExamplePullAll_2() {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	exclude := []string{"grayson", "tim"}

	result, err := PullAll(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "cassandra", "jason", "stephanie" }
}

func ExamplePullAt() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}

	result, err := PullAt(data, 1, 3)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 3.2, 5.2, 6.3 }
}

func ExampleRemove_1() {
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

func ExampleRemove_2() {
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

func ExampleReverse_1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	result, err := Reverse(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "tim", "grayson", "damian", "jason" }
}

func ExampleReverse_2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result, err := Reverse(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 9, 8, 7, 6, 5, 4, 3, 2, 1 }
}

func ExampleTail_1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	result, err := Tail(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "grayson", "tim" }
}

func ExampleTail_2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result, err := Tail(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 2, 3, 4, 5, 6, 7, 8, 9 }
}

func ExampleTake_1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	take := 2

	result, err := Take(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "jason", "damian" }
}

func ExampleTake_2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	take := 5

	result, err := Take(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 5 }
}

func ExampleTakeRight_1() {
	data := []string{"jason", "damian", "grayson", "tim"}
	take := 2

	result, err := TakeRight(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "grayson", "tim" }
}

func ExampleTakeRight_2() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	take := 5

	result, err := TakeRight(data, take)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 5, 6, 7, 8, 9 }
}

func ExampleUnion_1() {
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

func ExampleUnion_2() {
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

func ExampleUniq_1() {
	data := []string{"damian", "grayson", "grayson", "cassandra"}
	result, err := Uniq(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "grayson", "cassandra" }
}

func ExampleUniq_2() {
	data := []float64{1.1, 3.00000, 3.1, 2.2000000, 3, 2.2, 3.0}
	result, err := Uniq(data)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 3, 3.1, 2.2 }
}

func ExampleWithout_1() {
	data := []int{1, 2, 3, 4, 5, 6}
	exclude := []int{3}

	result, err := Without(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 4, 5, 6 }
}

func ExampleWithout_2() {
	data := []float64{1.1, 2.1, 3.2, 4.2, 5.2, 6.3}
	exclude := []float64{2.1, 3.2, 6.3}

	result, err := Without(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []float64{ 1.1, 4.2, 5.2 }
}

func ExampleWithout_3() {
	data := []string{"damian", "grayson", "cassandra", "tim", "tim", "jason", "stephanie"}
	exclude := []string{"grayson", "tim"}

	result, err := Without(data, exclude)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "cassandra", "jason", "stephanie" }
}
