package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Auth struct {
	Username string
	Password string
	IsAdmin  bool
}

type AuthStore struct {
	Users []Auth
}

func AuthCheck(user, password string ) (bool, bool){
	var aStore AuthStore
	aStore.LoadFromFile("users.json")

	for _, v := range aStore.Users {
		if user == v.Username {
			if password == v.Password {
				if v.IsAdmin == true {
					return true,true			//validatedUser,IsAdmin
				}
				return true,false			//ValidatedUser,IsNotAdmin
			}
		}
	}
	return false,false						//NotValidatedUser,IsNotAdmin
}

func (repo *AuthStore) LoadFromFile(file string) error {
	data, _ := ioutil.ReadFile(file)

	err := json.Unmarshal(data, &repo.Users)
	if err != nil {
		return errors.New("error in FileIO")
	}
	return nil
}
