package controllers

import (
	"bus-router-backend/mvc/services"
	"bus-router-backend/mvc/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func AddStation(w http.ResponseWriter, r *http.Request) {
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
	name, lat, lng, passengers := r.FormValue("name"), r.FormValue("lat"), r.FormValue("lng"), r.FormValue("passengers")
	if name == "" || lat == "" || lng == "" || passengers == "" { // via postman
		rBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			stationErr := &utils.AppErrors{
				Message:    "An error occured during get the credentials",
				StatusCode: http.StatusNoContent,
				Code:       "no content",
			}
			w.WriteHeader(stationErr.StatusCode)
			json.NewEncoder(w).Encode(stationErr)
			return
		}
		if station, err := services.ProcessJSONforStation(rBody); err != nil {
			stationErr := &utils.AppErrors{
				Message:    "An error occured during process for json to get data",
				StatusCode: http.StatusNoContent,
				Code:       "no content",
			}
			w.WriteHeader(stationErr.StatusCode)
			json.NewEncoder(w).Encode(stationErr)
			return
		} else {
			fmt.Println("my station --> ", station)
			if station.Name != "" && station.Latitude != 0 && station.Longitude != 0 {
				station, err := services.AddStation(station.Name, station.Latitude, station.Longitude, station.Passengers, services.MongoClient)
				if err != nil {
					fmt.Println("hata")
					w.WriteHeader(err.StatusCode)
					json.NewEncoder(w).Encode(err)
					return
				}
				fmt.Println(station)
				json.NewEncoder(w).Encode(station)
			} else {
				stationErr := &utils.AppErrors{
					Message:    "parse error occured",
					StatusCode: http.StatusBadGateway,
					Code:       "no content",
				}
				w.WriteHeader(stationErr.StatusCode)
				json.NewEncoder(w).Encode(stationErr)
				return
			}
		}
	}

}

// func DeleteStation(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseForm(); err != nil { // parse element from html
// 		adminErr := &utils.AppErrors{
// 			Message:    "An error occured during parsing the form",
// 			StatusCode: http.StatusBadGateway,
// 			Code:       "Parse error",
// 		}
// 		w.WriteHeader(adminErr.StatusCode)
// 		json.NewEncoder(w).Encode(adminErr)
// 		return
// 	}
// 	name := r.FormValue("name")
// 	if name == "" { // postman
// 		rBody, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			stationErr := &utils.AppErrors{
// 				Message:    "An error occured during get the credentials",
// 				StatusCode: http.StatusNoContent,
// 				Code:       "no content",
// 			}
// 			w.WriteHeader(stationErr.StatusCode)
// 			json.NewEncoder(w).Encode(stationErr)
// 			return
// 		}
// 		if station, err := services.ProcessJSONforStation(rBody); err != nil {
// 			stationErr := &utils.AppErrors{
// 				Message:    "An error occured during process for json to get data",
// 				StatusCode: http.StatusNoContent,
// 				Code:       "no content",
// 			}
// 			w.WriteHeader(stationErr.StatusCode)
// 			json.NewEncoder(w).Encode(stationErr)
// 			return
// 		} else {func DeleteStation(w http.ResponseWriter, r *http.Request)
// 			err := services.DeleteStation(station.Name, services.MongoClient)
// 			if err != nil {
// 				w.WriteHeader(err.StatusCode)
// 				json.NewEncoder(w).Encode(err)
// 				return
// 			}
// 			//fmt.Println(admin)
// 			json.NewEncoder(w).Encode("istasyon başarılı bir şekilde kaldırıldı \n")
// 		}
// 	}
// }
func DeleteStation(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		error := &utils.AppErrors{
			Message:    "url parse error occured",
			StatusCode: http.StatusBadGateway,
			Code:       "url parse erorr",
		}
		w.WriteHeader(error.StatusCode)
		json.NewEncoder(w).Encode(error)
		return
	}
	//fmt.Println(u)
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(m["name"][0])
	nameToDelete := m["name"][0]
	station, errs := services.DeleteStation(nameToDelete, services.MongoClient)
	fmt.Println("station to delete ", station)
	if errs != nil {
		w.WriteHeader(errs.StatusCode)
		json.NewEncoder(w).Encode(errs)
		return
	}
	json.NewEncoder(w).Encode("İstasyon başarılı bir şekilde kaldırıldı ! \n")
}

func UpdatePassengerOfStation(w http.ResponseWriter, r *http.Request) {
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
	passengers := r.FormValue("passengers")
	if passengers == "" {
		rBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			stationErr := &utils.AppErrors{
				Message:    "An error occured during get the credentials",
				StatusCode: http.StatusNoContent,
				Code:       "no content",
			}
			w.WriteHeader(stationErr.StatusCode)
			json.NewEncoder(w).Encode(stationErr)
			return
		}
		if station, err := services.ProcessJSONforStation(rBody); err != nil {
			stationErr := &utils.AppErrors{
				Message:    "An error occured during process for json to get data",
				StatusCode: http.StatusNoContent,
				Code:       "no content",
			}
			w.WriteHeader(stationErr.StatusCode)
			json.NewEncoder(w).Encode(stationErr)
			return
		} else {
			fmt.Println("station ---> ? ", station)
			if station.Name != "" && station.Passengers >= 0 {
				new_station, err := services.UpdatePassengerOfStation(station.Name, station.Passengers, services.MongoClient)
				if err != nil {
					w.WriteHeader(err.StatusCode)
					json.NewEncoder(w).Encode(err)
					return
				}
				fmt.Println(new_station)
				json.NewEncoder(w).Encode(new_station)
			} else {
				stationErr := &utils.AppErrors{
					Message:    "parse error occured",
					StatusCode: http.StatusBadGateway,
					Code:       "no content",
				}
				w.WriteHeader(stationErr.StatusCode)
				json.NewEncoder(w).Encode(stationErr)
				return
			}
		}
	}
}
