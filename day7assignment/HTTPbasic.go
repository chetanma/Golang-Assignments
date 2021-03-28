package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func HttpServerDemo() {
	hotelsRepo := HotelsRepo()

	http.HandleFunc("/hotel", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			idHotel, _ := strconv.Atoi(r.URL.Query().Get("id"))

			//update getHotel() will return error or not 
			//check if getHotel() return nil as error then marshal() wil handle it or not
			data, err := json.Marshal(hotelsRepo.GetHotel(idHotel))

			if err == nil{
				w.WriteHeader(http.StatusFound)
				fmt.Fprint(w, string(data))
			} else {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "error")
			}
			//TODO getHotelsById(id... int)
			//TODO getHotels()

		}
		if r.Method == "POST"{
			var hotel Hotel
            
            body,_ := ioutil.ReadAll(r.Body)
            json.Unmarshal(body, &hotel)
			fmt.Println(hotel)
            hotelsRepo.AddHotel(hotel)
            fmt.Println("before--->",hotelsRepo.GetHotels())

			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, "Hotel added successfully with id:",hotel.Id)

		}
		if r.Method == "PUT"{
			var hotel Hotel
			body,_ := ioutil.ReadAll(r.Body)
            json.Unmarshal(body, &hotel)
			hotelsRepo.UpdateHotel(hotel)
			fmt.Println("data--->",hotelsRepo.GetHotels())

			w.WriteHeader(http.StatusAccepted)
            fmt.Fprint(w, "Hotel updated successfully with id:",hotel.Id)
		}
		if r.Method =="DELETE"{
			HotelId, _ := strconv.Atoi(r.URL.Query().Get("id"))
			hotelsRepo.DeleteHotel(HotelId)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Hotel deleted successfully with id:",HotelId)
		}

	})

	http.ListenAndServe(":8000", nil)

}
