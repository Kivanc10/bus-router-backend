package controllers

import (
	"bus-router-backend/mvc/domain"
	"bus-router-backend/mvc/services"
	"bus-router-backend/mvc/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GetAdmin(w http.ResponseWriter, r *http.Request) {
	adminCode, err := strconv.ParseInt(r.URL.Query().Get("admin_code"), 10, 64)
	if err != nil {
		adminErr := &utils.AppErrors{
			Message:    "admin_code must be number",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(adminErr)
		return
	}
	admin, apiErr := services.GetAdmin(adminCode)
	if apiErr != nil {
		w.WriteHeader(apiErr.StatusCode)
		log.Printf(apiErr.Code)
		json.NewEncoder(w).Encode(apiErr)
		return
	}
	json.NewEncoder(w).Encode(admin)
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
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
	name, password := r.FormValue("name"), r.FormValue("password")

	if password == "" || name == "" { // via postman
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
		if admin, err := services.ProcessJSONforAdmin(rBody); err != nil {
			adminErr := &utils.AppErrors{
				Message:    "An error occured during process for json to get data",
				StatusCode: http.StatusNoContent,
				Code:       "no content",
			}
			w.WriteHeader(adminErr.StatusCode)
			json.NewEncoder(w).Encode(adminErr)
			return
		} else {
			admin, err := services.AdminLogin(admin)
			if err != nil {
				w.WriteHeader(err.StatusCode)
				json.NewEncoder(w).Encode(err)
				return
			}
			fmt.Println(admin)
			json.NewEncoder(w).Encode("Admin başarılı bir şekilde giriş yaptı ! \n")
		}
	}

}

func AdminInside(w http.ResponseWriter, r *http.Request) {
	admin := domain.Admin{
		Name:     os.Getenv("admin"),
		Password: os.Getenv("admin"),
	}
	err := services.AdminInside(admin)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Message)
		return
	}
	json.NewEncoder(w).Encode(admin)
}

func AdminLogout(w http.ResponseWriter, r *http.Request) {
	os.Setenv("admin", "")
	admin := domain.Admin{
		Name:     os.Getenv("admin"),
		Password: os.Getenv("admin"),
	}
	err := services.AdminLogout(admin)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode("Admin is logout succesfully\n")
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := services.GetAllUsers(services.MongoClient)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(users)
	return
}

func GetAllStations(w http.ResponseWriter, r *http.Request) {
	users, err := services.GetAllStations(services.MongoClient)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(users)
	return
}

func GetAllVoyages(w http.ResponseWriter, r *http.Request) {
	voyages, err := services.GetAllVoyages(services.MongoClient)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(voyages)
	return
}
