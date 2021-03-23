package main

import (
	"fmt"
	"strconv"
)

func main() {

	ch := make(chan Customer)

	go func() {
		defer close(ch)
		CustomerInitSilver(ch)
	}()
	out,ok:=<-ch
	for outTemp:= range out{
		if ok {
			fmt.Println(outTemp)
		}
	}
	fmt.Println()

}

//to initialize silver level customers only
func CustomerInitSilver(ch chan Customer) (<-chan Customer, error) {

	for i := 0; i < 16; i++ {
		cust, err := NewCustomer(("Cust" + strconv.Itoa(i)), (1500 + i*1), 120*i, 1)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		ch <- *cust
	}

	return ch, nil
}
