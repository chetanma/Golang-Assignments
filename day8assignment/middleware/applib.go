package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// type HotelData struct{
// 	Hotel <-chan []Hotel
// 	Err *error
// }

//function to get hotel by making external HTTP GET req
func FetchDataFromSourceUrl(url string) (<-chan []Hotel, error) {

	hotels :=make(chan []Hotel)
	 
	go func(){
		defer close(hotels)
		client := &http.Client{}
		payload := strings.NewReader(``)
		req, err := http.NewRequest("GET", url, payload)
		// if err!=nil{
		// 	return nilHotels, errors.New("bad request")		//new struct chan []Hotels,err
		// }
		req.Header.Add("Content-Type", "application/json")
		res, err := client.Do(req)
		if err!=nil{
			return 
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return 
		} else {
			hotelsTempData := make([]Hotel, 0, 100)
			json.Unmarshal(body, &hotelsTempData)
			hotels <- hotelsTempData
		}
	}()
	return hotels,nil
}

//function to book hotel by making external HTTP PUT req
func BookHotel(hotelByte []byte) (string,error) {

	var hotelJson Hotel
	json.Unmarshal(hotelByte, &hotelJson)
	payload := strings.NewReader(string(hotelByte))
	id := strconv.Itoa(hotelJson.Id)

	client := &http.Client{}
	url := hotelJson.Sourceurl + "/hotel/book?id=" + id
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return "",err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "",err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "",err
	}
	fmt.Println(string(body),string(strconv.Itoa(res.StatusCode)))

	return string(body),nil
}

//merging from two sources
func merge(chan1, chan2 <-chan []Hotel) <-chan []Hotel {
	ch := make(chan []Hotel)
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
