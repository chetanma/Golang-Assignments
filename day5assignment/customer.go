package main

import "errors"

type CustomerLevel int

const (
	Silver CustomerLevel = iota + 1
	Gold
	Platinum
)

type Customer struct {
	Name   string
	ID     int
	Points int
	Level  CustomerLevel
}

func NewCustomer(name string, id int, points int, level CustomerLevel) (*Customer, error) {

	if name == "" {
		return nil, errors.New("Empty name")
	}
	if id < 1000 {
		return nil, errors.New("ID value too low")
	}
	if points < 0 {
		return nil, errors.New("Points must be positive")
	}	
	return &Customer{
		Name:   name,
		ID:     id,
		Points: points,
		Level:  level,
	}, nil
}
