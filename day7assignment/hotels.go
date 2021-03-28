package main

import (
	"fmt"
	"strconv")


type HotelRepo interface {
	GetHotel(id int) *Hotel
	GetHotels() []Hotel
	AddHotel(hotel Hotel) string
	UpdateHotel(updatedHotel Hotel) string
	DeleteHotel(hotelId int) string
}

type HotelDataRepo struct {
	Data []Hotel
}

func HotelsRepo() HotelDataRepo {
	data := make([]Hotel, 0, 100)

	for i := 0; i < 10; i++ {
		data = append(data, Hotel{
			Name:     "Hotel" + strconv.Itoa(i),
			Id:       5000 + i,
			Capacity: 2 * i,
		})
	}

	return HotelDataRepo{
		Data: data,
	}
}

func (hotels *HotelDataRepo) GetHotel(id int) *Hotel {
	for _, hotel := range hotels.Data {
		if hotel.Id == id {
			return &hotel
		}
	}

	return nil
}

func (hotels *HotelDataRepo) GetHotels() []Hotel {
	return hotels.Data
}

func (hotelsRepo *HotelDataRepo) AddHotel(hotel Hotel) string {
	hotelsRepo.Data = append(hotelsRepo.Data, hotel)
	return "Hotel added successfully !!!"
}

func (hotelsRepo *HotelDataRepo) UpdateHotel(updatedHotel Hotel) string {
	for i,_ :=range hotelsRepo.Data{
		if hotelsRepo.Data[i].Id == updatedHotel.Id{
			hotelsRepo.Data[i] =updatedHotel
			//return true
		}
	}
	return "Hotel updated successfully !!!"			//return false if hotel not updated
}

func (hotelRepo *HotelDataRepo) DeleteHotel(hotelId int) string{
	updatedHotelRepo := make([]Hotel, 0, len(hotelRepo.Data)-1)
	
	for i,_ :=range hotelRepo.Data{
		if hotelRepo.Data[i].Id == hotelId{
			continue
		} else {
			updatedHotelRepo = append(updatedHotelRepo, hotelRepo.Data[i])

		}
	}
	fmt.Println("after-->",updatedHotelRepo)
	return "Hotel deleted successfully !!!"
}



