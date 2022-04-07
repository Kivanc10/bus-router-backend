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

func AddVoyage(w http.ResponseWriter, r *http.Request) {
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
	distance := r.FormValue("distance")
	if distance == "" {
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
		if voyage, err := services.ProcessJSONforVoyage(rBody); err != nil {
			voyageErr := &utils.AppErrors{
				Message:    "An error occured during process for json to get data",
				StatusCode: http.StatusNoContent,
				Code:       "no content",
			}
			w.WriteHeader(voyageErr.StatusCode)
			json.NewEncoder(w).Encode(voyageErr)
			return
		} else {
			if voyage.From.Name != "" && voyage.To.Name != "" && voyage.Distance != 0 {
				voyage, err := services.AddVoyage(voyage.From.Name, voyage.To.Name, voyage.Distance, services.MongoClient)
				if err != nil {
					fmt.Println("hata")
					w.WriteHeader(err.StatusCode)
					json.NewEncoder(w).Encode(err)
					return
				}
				fmt.Println(voyage)
				json.NewEncoder(w).Encode(voyage)
			} else {
				voyageErr := &utils.AppErrors{
					Message:    "parse error occured",
					StatusCode: http.StatusBadGateway,
					Code:       "no content",
				}
				w.WriteHeader(voyageErr.StatusCode)
				json.NewEncoder(w).Encode(voyageErr)
				return
			}
		}
	}
}

func DeleteVoyage(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(m)
	if len(m) == 2 {
		from, to := m["from"][0], m["to"][0]
		fmt.Println(from, to)
		voyage, err := services.DeleteVoyage(from, to, services.MongoClient)
		fmt.Println("voyage to delete ", voyage)

		if err != nil {
			w.WriteHeader(err.StatusCode)
			json.NewEncoder(w).Encode(err)
			return
		}
		json.NewEncoder(w).Encode("Sefer başarılı bir şekilde kaldırıldı ! \n")

	}
}
