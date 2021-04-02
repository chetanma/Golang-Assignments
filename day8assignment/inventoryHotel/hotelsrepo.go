package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
)

type HotelStoreJson struct {
	mutex sync.RWMutex
}

type HotelRepo interface {
	GetHotel(id int) ([]Hotel, error)
	GetHotels() ([]Hotel, error)
	AddHotel(hotel Hotel) error
	UpdateHotel(updatedHotel Hotel) (bool, error)
	DeleteHotel(hotelId int) (bool, error)
}

type HotelDataRepo struct {
	Data []Hotel
}

//get all hotels list
func (store *HotelStoreJson) GetHotels() ([]Hotel, error) {

	store.mutex.RLock()
	defer store.mutex.RUnlock()

	data, _ := ioutil.ReadFile("data.json")
	var hotels []Hotel //if hotels is taken as type of HotelDataRepo, data is getting null
	json.Unmarshal(data, &hotels)
	return hotels, nil
}

//get a hotel
func (store *HotelStoreJson) GetHotel(id int) ([]Hotel, error) {
	store.mutex.RLock()
	defer store.mutex.RLock()

	//transfer this readFile content to LoadFromFile()
	data, _ := ioutil.ReadFile("data.json")
	var hotels []Hotel
	json.Unmarshal(data, &hotels)

	for _, hotel := range hotels {
		if hotel.Id == id {
			hotelFound := make([]Hotel, 0, 100)
			hotelFound = append(hotelFound, hotel)
			return hotelFound, nil
		}
	}
	return nil, errors.New("Hotel not found")
}

//add a hotel
func (store *HotelStoreJson) AddHotel(hotel Hotel) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	data, _ := ioutil.ReadFile("data.json")
	var hotels []Hotel
	json.Unmarshal(data, &hotels)
	hotels = append(hotels, hotel)
	data, _ = json.Marshal(hotels)
	err := ioutil.WriteFile("data.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

//update a hotel
func (store *HotelStoreJson) UpdateHotel(hotelId int, updatedHotel Hotel) (bool, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	var found bool
	data, _ := ioutil.ReadFile("data.json")
	var hotels []Hotel
	json.Unmarshal(data, &hotels)
	for i, _ := range hotels {
		if hotels[i].Id == hotelId {
			found = true
			hotels[i] = updatedHotel
		}
	}
	if !found {
		return found, errors.New("hotel not found with provided ID")
	} else {
		//transfer this repetative code to WriteToFile
		data, _ = json.Marshal(hotels)
		err := ioutil.WriteFile("data.json", data, 0644)
		if err != nil {
			return true, errors.New("error in FileIO process")
		} else {
			return true, nil
		}
	}
}

//book a hotel
func (store *HotelStoreJson) BookHotel(hotelId int) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	var found bool

	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		return errors.New("error in FileIO process")
	}
	var hotels []Hotel
	json.Unmarshal(data, &hotels)
	for i, _ := range hotels {
		if hotels[i].Id == hotelId {
			found = true
			if hotels[i].Availability == 0 {
				return errors.New("not available")
			}
			hotels[i].Availability -= 1
		}
	}
	if found == false {
		return errors.New("not found")
	}
	//transfer this redundent code to WriteToFile
	data, _ = json.Marshal(hotels)
	err1 := ioutil.WriteFile("data.json", data, 0644)
	if err1 != nil {
		return errors.New("error in FileIO process")
	}
	return nil

}

//delete a hotel
func (store *HotelStoreJson) DeleteHotel(hotelId int) (bool, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	data, _ := ioutil.ReadFile("data.json")
	var hotels []Hotel
	json.Unmarshal(data, &hotels)

	UpdatedHotels := make([]Hotel, 0, 100)
	var found bool
	for i, _ := range hotels {
		if hotels[i].Id == hotelId {
			found = true
			continue
		} else {
			UpdatedHotels = append(UpdatedHotels, hotels[i])
		}
	}

	if !found {
		return false, errors.New("hotel not found with provided ID")
	} else {
		data, _ = json.Marshal(UpdatedHotels)
		err := ioutil.WriteFile("data.json", data, 0644)
		if err != nil {
			return true, errors.New("error in FileIO process")
		} else {
			return true, nil
		}
	}
}
