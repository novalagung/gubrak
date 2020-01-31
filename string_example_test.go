package gubrak

import (
	"fmt"
)

func ExampleRandomString() {
	result := RandomString(32)
	fmt.Println(result) // generates random 32 character like: YodQeljldGFluOhaHrlWdICKDtDHSvzA
}

func ExampleReplaceCaseInsensitive() {
	result := ReplaceCaseInsensitive("lOrEm IPsUm DoLor Sit AMEt", "ipsum", "batman")
	fmt.Println(result) // lOrEm batman DoLor Sit AMEt
}
