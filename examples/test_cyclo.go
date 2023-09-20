package examples

import (
	"fmt"
)

func C(a byte, b int, c string, d uint8, e []byte, f uintptr, g bool, h rune, i complex64, j float32) {
	if b < 0 {
		fmt.Println("a")
	} else if b == 5 {
		fmt.Println("b")
	} else if b == 6 {
		fmt.Println("c")
	} else {
		fmt.Println("d")
	}
	fmt.Println(b)


	switch b {
	case 1:
		fmt.Println("a")
		fmt.Println("a")
	case 2:
		fmt.Println("b")
	case 3:
		fmt.Println("c")
	default:
		fmt.Println("d")
	}
}
