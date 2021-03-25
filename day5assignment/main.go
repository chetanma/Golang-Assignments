package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type LevelUpgrade func(*Customer) bool
type BonusPoints func(*Customer) bool

func main() {

	CustomerPipeline()

}

//demonstration of Pipeline
func CustomerPipeline() {
	ch := Generator()

	inA1 := UpgradeLevelWorker(ch, LevelUpgradeClosure())
	inA2 := UpgradeLevelWorker(ch, LevelUpgradeClosure())

	outA := merge(inA1, inA2)

	inB1 := BonusPointsWorker(outA, BounusPointsclosure())
	inB2 := BonusPointsWorker(outA, BounusPointsclosure())

	outB := merge(inB1, inB2)

	
	WriteOutputToJson(outB,"file.json")
}

//to initialize silver level customers only
func Generator() <-chan Customer {
	ch := make(chan Customer)
	go func() {
		defer close(ch)
		for i := 0; i < 16; i++ {
			cust, err := NewCustomer(("Cust" + strconv.Itoa(i)), (1500 + i*1), 120*i, 1)
			if err != nil {
				fmt.Println(err)
				continue
			}
			ch <- *cust
		}
	}()
	return ch
}

//Level upgrade closure
func LevelUpgradeClosure() LevelUpgrade {
	return func(c *Customer) bool {
		if (800 < c.Points) && (c.Points < 1500) {
			c.Level = Gold
			return true
		} else if c.Points >= 1500 {
			c.Level = Platinum
			return true
		}
		return false
	}
}

//will upgrade customer level
func UpgradeLevelWorker(in <-chan Customer, fun LevelUpgrade) <-chan Customer {
	ch := make(chan Customer)

	go func() {
		defer close(ch)
		for out := range in {
			if fun(&out) {
				ch <- out
			}
		}
	}()

	return ch
}

//bounus points closure
func BounusPointsclosure() BonusPoints {
	return func(c *Customer) bool {
		if c.Level == Gold {
			c.Points += c.Points * 10 / 100
			return true
		}
		if c.Level == Platinum {
			c.Points += c.Points * 25 / 100
			return true
		}
		return false
	}
}

//assign bonus points
func BonusPointsWorker(in <-chan Customer, fun BonusPoints) <-chan Customer {
	ch := make(chan Customer)

	go func() {
		defer close(ch)
		for out := range in {
			if fun(&out) {
				ch <- out
			}
		}
	}()
	return ch
}

//Merging channels to a single channel--->Fan-in pattern
func merge(chan1, chan2 <-chan Customer) <-chan Customer {
	ch := make(chan Customer)
	go func() {
		defer close(ch)
		for chan1 != nil || chan2 != nil {
			select {
			case v, ok := <-chan1:
				if ok {
					ch <- v
				} else {
					chan1 = nil
				}
			case v, ok := <-chan2:
				if ok {
					ch <- v
				} else {
					chan2 = nil
				}
			}
		}
	}()
	
	return ch
}

func  WriteOutputToJson(out <-chan Customer, file string){
	CustomerSlice := make([]Customer, 0,100)
	for temp:=range out{
		CustomerSlice = append(CustomerSlice, temp)
		fmt.Println("cust->",temp)
	}
	
	data, _ := json.Marshal(CustomerSlice)
	
	fmt.Println("json data-->",string(data))

	ioutil.WriteFile(file, data, 0644)

	return
}