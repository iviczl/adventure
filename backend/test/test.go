package test

import (
	"fmt"
	"os"
)

func test() {
	fmt.Println(os.Args[0])
	if len(os.Args) > 1 {
		for index, arg := range os.Args[1:] {
			fmt.Println(index, arg)
		}
	} else {
		var currency string
		for len(currency) == 0 {
			fmt.Print("Enter the currency code:")
			fmt.Scan(&currency)
			fmt.Println(currency)
		}
	}
	fmt.Printf("end")
}
