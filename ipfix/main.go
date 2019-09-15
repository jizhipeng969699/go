package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Unix(1111111111111111111,0).UTC())

	sli:=[]int{1,2,3,4,5,6,7,8,9,10}

	sli = sli[8:]

	fmt.Println(sli)

	lentest(sli)
}

func lentest(sli []int) {
	fmt.Println(len(sli))
}

func test(sli []int) {
	for k, v := range sli {
		fmt.Println(k,v)
	}

	fmt.Println(sli[0],sli[1])
}
