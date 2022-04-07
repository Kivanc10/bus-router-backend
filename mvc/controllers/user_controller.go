package controllers

import (
	"bus-router-backend/mvc/services"
	"bus-router-backend/mvc/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { // parse element from html
		adminErr := &utils.AppErrors{
			Message:    "An error occured during parsing the form",
			StatusCode: http.StatusBadGateway,
			Code:       "Parse error",
		}
		w.WriteHeader(adminErr.StatusCode)
		json.NewEncoder(w).Encode(adminErr)
		return
	}
	name, email, password := r.FormValue("name"), r.FormValue("email"), r.FormValue("password")
	if name == "" || email == "" || password == "" { // via postman
		rBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			adminErr := &utils.AppErrors{
				Message:    "An error occured during get the credentials",
				StatusCode: http.StatusNoContent,
				Code:       "no content",
			}
			w.WriteHeader(adminErr.StatusCode)
			json.NewEncoder(w).Encode(adminErr)
			return
		}
		if user, err := services.ProcessJSONforUser(rBody); err != nil {
			adminErr := &utils.AppErrors{
				Message:    "An error occured during process for json to get data",
				StatusCode: http.StatusNoContent,
				Code:       "no content",
			}
			w.WriteHeader(adminErr.StatusCode)
			json.NewEncoder(w).Encode(adminErr)
			return
		} else {
			if user.Name != "" && user.Email != "" && user.Password != "" {
				user, err := services.SignUp(user.Name, user.Email, user.Password, services.MongoClient)
				if err != nil {
					w.WriteHeader(err.StatusCode)
					json.NewEncoder(w).Encode(err)
					return
				}
				fmt.Println(user)
				json.NewEncoder(w).Encode(user)
			}
		}

	}

}

func SignIn(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { // parse element from html
		adminErr := &utils.AppErrors{
			Message:    "An error occured during parsing the form",
			StatusCode: http.StatusBadGateway,
			Code:       "Parse error",
		}
		w.WriteHeader(adminErr.StatusCode)
		json.NewEncoder(w).Encode(adminErr)
		return
	}
	password, name := r.FormValue("password"), r.FormValue("name")
	if password == "" || name == "" {
		rBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			adminErr := &utils.AppErrors{
				Message:    "An error occured during get the credentials",
				StatusCode: http.StatusNoContent,
				Code:       "no content",
			}
			w.WriteHeader(adminErr.StatusCode)
			json.NewEncoder(w).Encode(adminErr)
			return
		}
		//
		if user, err := services.ProcessJSONforUser(rBody); err != nil {
			adminErr := &utils.AppErrors{
				Message:    "An error occured during process for json to get data",
				StatusCode: http.StatusNoContent,
				Code:       "no content",
			}
			w.WriteHeader(adminErr.StatusCode)
			json.NewEncoder(w).Encode(adminErr)
			return
		} else {
			if user.Name != "" && user.Password != "" {
				user, err := services.SignIn(user.Name, user.Password, services.MongoClient)
				if err != nil {
					w.WriteHeader(err.StatusCode)
					json.NewEncoder(w).Encode(err)
					return
				}
				fmt.Println(user)
				json.NewEncoder(w).Encode(user)
			}
		}
	}
}

func GetInside(w http.ResponseWriter, r *http.Request) {
	user, err := services.GetInside(services.MongoClient)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(user)
	return
}

func LogoutForUser(w http.ResponseWriter, r *http.Request) {
	os.Setenv("Token", "")
	os.Setenv("username", "")
	err := services.LogoutForUser()
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode("user is logout succesfully\n")
}
