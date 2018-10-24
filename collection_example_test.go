package gubrak

import (
	"fmt"
	"strings"
)

func ExampleCount_countMap1() {
	data := map[string]interface{}{
		"name":   "jason",
		"age":    12,
		"isMale": true,
	}
	result, err := Count(data)
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
	// ===> 1
}

func ExampleCount_countSlice1() {
	data := []string{"damian", "grayson", "cassandra"}
	result, err := Count(data)
	// ===> 3
}

func ExampleCount_countSlice2() {
	data := []string{"damian", "grayson", "cassandra"}

	result, err := Count(data, func(each string) bool {
		return strings.Contains(each, "d")
	})
	// ===> 2
}

func ExampleCount_countSlice3() {
	data := []string{"damian", "grayson", "cassandra"}

	result, err := Count(data, func(each string, i int) bool {
		return len(each) > 6 && i > 1
	})
	// ===> 1
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
}

func ExampleEach_eachSlice1() {
	data := []string{"damian", "grayson", "cassandra"}

	err := Each(data, func(each string) {
		fmt.Println(each)
	})
}

func ExampleEach_eachSlice2() {
	data := []string{"damian", "grayson", "cassandra"}

	err := Each(data, func(each string, i nt) {
		fmt.Printf("element %d: %s \n", i, each)
	})
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
}

func ExampleEachRight_eachRightSlice1() {
	data := []string{"damian", "grayson", "cassandra"}

	err := EachRight(data, func(each string) {
		fmt.Println(each)
	})
}

func ExampleEachRight_eachRightSlice2() {
	data := []string{"damian", "grayson", "cassandra"}

	err := EachRight(data, func(each string, i nt) {
		fmt.Printf("element %d: %s \n", i, each)
	})
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

	// ===> Sample { EbookName: "rework", DailyDownloads: 12000 }
}

func ExampleFind_find2() {
	data := []string{"clean code", "rework", "detective comics"}

	result, err := Find(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 1)

	// ===> "detective comics"
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
		return strings.Contains(each, "co")
	})

	// ===> Sample { EbookName: "detective comics", DailyDownloads: 11500 }
}

func ExampleFindLast_findLast2() {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	result, err := FindLast(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 2)

	// ===> "detective comics"
}

func ExampleFindLast_findLast3() {
	data := []string{"clean code", "rework", "detective comics", "coco"}

	result, err := FindLast(data, func(each string, i int) bool {
		return strings.Contains(each, "co")
	}, 3)

	// ===> "coco"
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
	// ===> true
}

func ExampleIncludes_includesMap2() {
	data := map[string]string{
		"name":  "grayson",
		"hobby": "helping people",
	}

	result, err := Includes(data, "batmobile")
	// ===> false
}

func ExampleIncludes_includesSlice1() {
	data := []string{"damian", "tim", "jason", "grayson"}
	result, err := Includes(data, "tim")
	// ===> true
}

func ExampleIncludes_includesSlice2() {
	data := []string{"damian", "tim", "jason", "grayson"}
	result, err := Includes(data, "tim", 2)
	// ===> false
}

func ExampleIncludes_includesSlice3() {
	data := []string{"damian", "tim", "jason", "grayson"}
	result, err := Includes(data, "cassandra")
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
