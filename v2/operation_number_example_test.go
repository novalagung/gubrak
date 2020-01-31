package gubrak

import "fmt"

func ExampleRandomInt() {
	result := RandomInt(10, 12)
	fmt.Println(result) // generates random int between 10 to 12 like: 10 or 11 or 12
}
