package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Fetch Hotels Handler
func FetchHotelsHandler(sources []string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		hotelList := make([]Hotel, 0, 100)
		mapValues := r.URL.Query()
		hotelIdMap := mapValues["id"]
		sizeMap := mapValues["size"]

		//if no params passed, get all Hotels...
		if hotelIdMap == nil {

			//source 1 -> on port 8001
			url1 := sources[0] + "/hotel"
			hotelListFromSource1, err := FetchDataFromSourceUrl(url1)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(rw, "internal server error")
				return
			}
			//source 2 -> on port 8002
			url2 := sources[1] + "/hotel"
			hotelListFromSource2, err := FetchDataFromSourceUrl(url2)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(rw, "internal server error")
				return
			}

			//take above code in a single for loop, append channel to the slice which will have channel, i.e.,
			// it will be slice of channels and then pass it to the merge as below:
			//-------merge(hotelChanSlice...)------

			//merging data from both sources
			hotelListChan := merge(hotelListFromSource1, hotelListFromSource2)
			for hotel := range hotelListChan {
				for i, _ := range hotel {
					hotelList = append(hotelList, hotel[i])
				}
			}
		}
		// else {
		// 	//get only hotel with provided id
		// 	HotelId := mapValues.Get("id")
		// 	for _, source := range sources {
		// 		url := source + "/hotel?id=" + HotelId
		// 		hotelFromSource, err := FetchDataFromSourceUrl(url)
		// 		if err != nil {
		// 			rw.WriteHeader(http.StatusInternalServerError)
		// 			fmt.Fprint(rw, "internal server error")
		// 			return
		// 		} else {
		// 			hotelList = append(hotelList, hotelFromSource...)
		// 		}
		// 	}
		// }

		//Size Filter-->parameters allowed -> S, L, XL
		if sizeMap != nil {
			hotelListFilteredBySize := make([]Hotel, 0, 100)
			for _, hotel := range hotelList {
				for _, size := range sizeMap {
					if hotel.Size == size {
						hotelListFilteredBySize = append(hotelListFilteredBySize, hotel)
					}
				}
			}
			json.NewEncoder(rw).Encode(hotelListFilteredBySize)
			hotelListByte, _ := json.Marshal(hotelListFilteredBySize)
			rw.WriteHeader(http.StatusOK)
			fmt.Fprint(rw, string(hotelListByte))
			return
		}

		hotelListByte, _ := json.Marshal(hotelList)
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, string(hotelListByte))
	}
}

//Book Hotel Handler
func BookHotelHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		body, _ := ioutil.ReadAll(r.Body)
		res, err := BookHotel(body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		switch res {
		case "not available":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Unable to book the room. Room is not available !!")
			return
		case "not found":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Hotel not found with provided ID !!")
			return
		case "booking successful":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Hotel booked successfully !!")
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

	}
}
