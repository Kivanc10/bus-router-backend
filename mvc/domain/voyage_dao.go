package domain

import (
	"bus-router-backend/mvc/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProcessJSONforVoyage(body []byte) (*Voyage, error) {
	var voyage Voyage
	var temp map[string]interface{}
	if err := json.Unmarshal(body, &temp); err != nil {
		return &Voyage{}, err
	}
	if temp["from"] != nil {
		voyage.From.Name = temp["from"].(string)
	}
	if temp["to"] != nil {
		voyage.To.Name = temp["to"].(string)
	}
	if temp["distance"] != nil {
		if s, ok := temp["distance"].(float64); ok {
			voyage.Distance = int64(s)
		}
	}
	return &voyage, nil
}

func insertVoyage(client *mongo.Client, voyage *Voyage) *utils.AppErrors {
	collection := client.Database("User_passenger").Collection("voyages")
	toSave, err := bson.Marshal(voyage)
	if err != nil {
		stationError := &utils.AppErrors{
			Message:    "unable to marshall the station struct",
			StatusCode: http.StatusBadRequest,
			Code:       "marshall error",
		}
		return stationError
	}
	n := []int{1, 2}
	for i, _ := range n {
		result := false
		if i == 0 {
			_, result = isAlreadyExistVoyage(voyage.From, voyage.To, client)
			if !result {
				res, err := collection.InsertOne(context.Background(), toSave)
				if err != nil {
					stationError := &utils.AppErrors{
						Message:    "unable to insert voyage",
						StatusCode: http.StatusBadRequest,
						Code:       "mongo insert error",
					}
					return stationError
				}
				id := res.InsertedID
				fmt.Println("id --> ", id)
			} else {
				return &utils.AppErrors{
					Message:    "the voyage is already exist",
					StatusCode: http.StatusUnauthorized,
					Code:       "already exist voyage is provided",
				}
			}
			//fmt.Println(result)
		} else {
			_, result = isAlreadyExistVoyage(voyage.To, voyage.From, client)
			if !result {
				newVoyage := *voyage
				newVoyage.From = voyage.To
				newVoyage.To = voyage.From
				toSave, err := bson.Marshal(newVoyage)
				if err != nil {
					stationError := &utils.AppErrors{
						Message:    "unable to marshall the station struct",
						StatusCode: http.StatusBadRequest,
						Code:       "marshall error",
					}
					return stationError
				}
				res, err := collection.InsertOne(context.Background(), toSave)
				if err != nil {
					stationError := &utils.AppErrors{
						Message:    "unable to insert voyage",
						StatusCode: http.StatusBadRequest,
						Code:       "mongo insert error",
					}
					return stationError
				}
				id := res.InsertedID
				fmt.Println("id --> ", id)
			}

		}
	}
	return nil
}

// get stations by name and collect them into voyage
func AddVoyage(from, to string, distance int64, client *mongo.Client) (*Voyage, *utils.AppErrors) {
	fromStation, err := GetStationByName(from, client)
	if err != nil {
		stationErr := &utils.AppErrors{
			Message:    "there is no station exist that has the name " + from,
			StatusCode: http.StatusNotFound,
			Code:       "no station exist",
		}
		return &Voyage{}, stationErr
	}
	toStation, err := GetStationByName(to, client)
	if err != nil {
		stationErr := &utils.AppErrors{
			Message:    "there is no station exist that has the name " + to,
			StatusCode: http.StatusNotFound,
			Code:       "no station exist",
		}
		return &Voyage{}, stationErr
	}
	var voyage Voyage
	voyage.From = *fromStation
	voyage.To = *toStation
	voyage.Distance = distance
	if err := insertVoyage(client, &voyage); err != nil {
		return &Voyage{}, err
	}
	return &voyage, nil
}

func isAlreadyExistVoyage(from, to Station, client *mongo.Client) (Voyage, bool) {
	collection := client.Database("User_passenger").Collection("voyages")
	temp := bson.M{"from": bson.M{"$eq": from}, "to": bson.M{"$eq": to}}
	result := Voyage{}
	err := collection.FindOne(context.Background(), temp).Decode(&result)
	if err != nil {
		return Voyage{}, false
	}
	return result, true

}

// add two voyages at once --> as = from to, to from

func DeleteVoyage(from, to string, client *mongo.Client) (Voyage, *utils.AppErrors) {
	collection := client.Database("User_passenger").Collection("voyages")
	from_, err := GetStationByName(from, client)
	fmt.Println("from --> ", from_)
	if err != nil {
		return Voyage{}, &utils.AppErrors{
			Message:    "an error occured during delete the voyage",
			StatusCode: http.StatusBadRequest,
			Code:       "unable to delete voyage",
		}
	}
	to_, err := GetStationByName(to, client)
	fmt.Println("to --> ", to_)
	if err != nil {
		return Voyage{}, &utils.AppErrors{
			Message:    "an error occured during delete the voyage",
			StatusCode: http.StatusBadRequest,
			Code:       "unable to delete voyage",
		}
	}
	temp := bson.M{"from": bson.M{"$eq": from_}, "to": bson.M{"$eq": to_}}
	fmt.Println("voyage to del ----> ", temp)
	var voyage Voyage
	new_err := collection.FindOneAndDelete(context.Background(), temp).Decode(&voyage)
	if new_err != nil {
		return Voyage{}, &utils.AppErrors{
			Message:    "an error occured during delete the stations",
			StatusCode: http.StatusBadRequest,
			Code:       "unable to delete station",
		}
	}
	return voyage, nil
}
