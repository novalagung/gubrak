package gubrak

import (
	"fmt"
	"log"
)

func ExampleChunk() {
	data := []int{1, 2, 3, 4, 5}
	size := 2

	result, err := Chunk(data, size)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> [][]int{ { 1, 2 }, { 3, 4 }, { 5 } }
}

func ExampleCompact() {
	data := []int{-2, -1, 0, 1, 2}

	result, err := Compact(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ -2, -1, 1, 2 }
}

func ExampleConcat() {
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

func ExampleDifference() {
	data := []int{1, 2, 3, 4, 4, 6, 7}
	dataDiff := []int{2, 7}

	result, err := Difference(data, dataDiff)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 3, 4, 4, 6 }
}

func ExampleDrop() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	n := 1

	result, err := Drop(data, n)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 2, 3, 4, 4, 5, 6 }
}

func ExampleDropRight() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	n := 1

	result, err := DropRight(data, n)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 1, 2, 3, 4, 4, 5 }
}

func ExampleFill() {
	data := []int{1, 2, 3, 4, 4, 5, 6}
	replacement := 9

	result, err := Fill(data, replacement)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []int{ 9, 9, 9, 9, 9, 9, 9 }
}

func ExampleFindIndex() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}
	result, err := FindIndex(data, func(each string) bool {
		return each == "tim"
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> 3
}

func ExampleFindLastIndex() {
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

func ExampleFromPairs() {
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

func ExampleFirst() {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := First(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> "damian"
}

func ExampleIndexOf() {
	data := []string{"damian", "grayson", "cass", "tim", "tim", "jason", "steph"}

	IndexOf(data, "duke")    // ===> -1
	IndexOf(data, "tim")     // ===> 3
	IndexOf(data, "tim", 4)  // ===> 4
	IndexOf(data, "tim", -4) // ===> 3
	IndexOf(data, "tim", -3) // ===> 4
	IndexOf(data, "tim", -2) // ===> -1
}

func ExampleInitial() {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Initial(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(result)
	// ===> []string{ "damian", "grayson" }
}
