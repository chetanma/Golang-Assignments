package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

//Fetch Hotel Handler function
func FetchHotelsHandler(store *HotelStoreJson) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		mapValues := r.URL.Query()

		hotelIdMap := mapValues["id"]

		//if no params passed, get all Hotels...
		if hotelIdMap == nil {
			hotels, err := store.GetHotels()
			hotelBytes, _ := json.Marshal(hotels)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "No hotels found !!!")
			} else {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string(hotelBytes))
			}
		} else {
			//get only hotel with provided id
			HotelId, _ := strconv.Atoi(mapValues.Get("id"))
			hotel, err := store.GetHotel(HotelId)
			hotelBytes, _ := json.Marshal(hotel)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "No hotel found with given ID!!!")
			} else {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string(hotelBytes))
			}
		}
	}
}

//Create Hotel Handler function
func CreateHotelsHandler(store *HotelStoreJson) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		isAdminCtx := ctx.Value("isAdmin")
		isAdmin := isAdminCtx.(bool)
		if !isAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Not Authorized")
			return
		}

		w.Header().Set("Content-Type", "application/json")

		var hotel Hotel
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &hotel)
		err := store.AddHotel(hotel)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		} else {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, "Record saved successfully!")
		}
	}
}

//Edit Hotel Handler function
func EditHotelsHandler(store *HotelStoreJson) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		isAdminCtx := ctx.Value("isAdmin")
		isAdmin := isAdminCtx.(bool)
		if !isAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Not Authorized")
			return
		}

		w.Header().Set("Content-Type", "application/json")

		HotelIdString := r.URL.Query().Get("id") //skip
		var UpdatedHotel Hotel
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &UpdatedHotel)
		HotelId, _ := strconv.Atoi(HotelIdString)

		found, err := store.UpdateHotel(HotelId, UpdatedHotel)

		if found == false {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, err)
		} else {
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
			} else {
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprint(w, "Record edited successfully !")
			}
		}
	}
}

//Book Hotel Handler function
func BookHotelHandler(store *HotelStoreJson) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		HotelIdString := r.URL.Query().Get("id")
		HotelId, _ := strconv.Atoi(HotelIdString)

		 err := store.BookHotel(HotelId)

		// if found == false {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	fmt.Fprint(w, "Hotel not found with provided ID")
		// 	return
		// }
		if err != nil {
			switch err.Error() {
			case "not available":
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "not available")	//Unable to book the room. Room are not available
				return
			case "not found":
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "not found")	//Hotel not found with provided ID
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
				return
			}
		}
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, "booking successful")
	}
}

//Delete Hotel Handler function
func DeleteHotelsHandler(store *HotelStoreJson) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		isAdminCtx := ctx.Value("isAdmin")
		isAdmin := isAdminCtx.(bool)
		if !isAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Not Authorized")
			return
		}

		w.Header().Set("Content-Type", "application/json")

		HotelIdString := r.URL.Query().Get("id")
		HotelId, _ := strconv.Atoi(HotelIdString)

		found, err := store.DeleteHotel(HotelId)

		if found == false {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, err)
		} else {
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
			} else {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "Record deleted successfully!")
			}
		}
	}
}

//Basic Auth Middleware
func BasicAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		verifiedUser, isAdmin := AuthCheck(user, pass)
		if !ok || !verifiedUser {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Not Authorized User")
			return
		}
		context := context.WithValue(r.Context(), "isAdmin", isAdmin)
		h(w, r.WithContext(context))
	}
}
