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

func insertStation(client *mongo.Client, station *Station) *utils.AppErrors {
	collection := client.Database("User_passenger").Collection("stations")
	toSave, err := bson.Marshal(station)
	if err != nil {
		stationError := &utils.AppErrors{
			Message:    "unable to marshall the station struct",
			StatusCode: http.StatusBadRequest,
			Code:       "marshall error",
		}
		return stationError
	}
	if _, result := isAlreadyExistStation(station.Name, client); !result {
		res, err := collection.InsertOne(context.Background(), toSave)
		if err != nil {
			stationError := &utils.AppErrors{
				Message:    "unable to insert station",
				StatusCode: http.StatusBadRequest,
				Code:       "mongo insert error",
			}
			return stationError
		}
		id := res.InsertedID
		fmt.Println("id --> ", id)

	} else {
		return &utils.AppErrors{
			Message:    "the station is already exist",
			StatusCode: http.StatusUnauthorized,
			Code:       "already exist station is provided",
		}
	}
	return nil
}

func isAlreadyExistStation(name string, client *mongo.Client) (Station, bool) {
	collection := client.Database("User_passenger").Collection("stations")
	temp := bson.M{"name": bson.M{"$eq": name}}
	result := Station{}
	err := collection.FindOne(context.Background(), temp).Decode(&result)
	if err != nil {
		return Station{}, false
	}
	return result, true
}

func AddStation(name string, lat, lng float64, passengers int64, client *mongo.Client) (*Station, *utils.AppErrors) {
	var station Station
	station.Name = name
	station.Latitude = lat
	station.Longitude = lng
	station.Passengers = passengers
	if err := insertStation(client, &station); err != nil {
		return &Station{}, err
	}
	return &station, nil
}

func DeleteStation(name string, client *mongo.Client) (Station, *utils.AppErrors) {
	collection := client.Database("User_passenger").Collection("stations")
	temp := bson.M{"name": bson.M{"$eq": name}}
	var station Station
	err := collection.FindOneAndDelete(context.Background(), temp).Decode(&station)
	fmt.Println("err -> ", err)
	if err != nil {
		return Station{}, &utils.AppErrors{
			Message:    "an error occured during delete the stations",
			StatusCode: http.StatusBadRequest,
			Code:       "unable to delete station",
		}
	}
	return station, nil
}

func ProcessJSONforStation(body []byte) (*Station, error) {
	fmt.Println("provided body --> ", string(body))
	var station Station
	var temp map[string]interface{}
	if err := json.Unmarshal(body, &temp); err != nil {
		return &Station{}, err
	}
	fmt.Println("in station --> ", temp)
	if temp["name"] != nil {
		station.Name = temp["name"].(string)
	}
	if temp["lat"] != nil {
		station.Latitude = temp["lat"].(float64)
	}
	if temp["lng"] != nil {
		station.Longitude = temp["lng"].(float64)
	}
	if temp["passengers"] != nil {
		if s, ok := temp["passengers"].(float64); ok {
			fmt.Println(s, ok)
			station.Passengers = int64(s)
			fmt.Println("station.pass -_> ", station.Passengers)
		}
	}

	return &station, nil
}

func UpdatePassengerOfStation(name string, passengers int64, client *mongo.Client) (*Station, *utils.AppErrors) {
	collection := client.Database("User_passenger").Collection("stations")
	temp := bson.M{"name": bson.M{"$eq": name}}
	res := collection.FindOne(context.Background(), temp)
	var station Station
	if err := res.Decode(&station); err != nil {
		stationErr := &utils.AppErrors{
			Message:    "an error occured during decode the station to update",
			StatusCode: http.StatusBadGateway,
			Code:       "decode error",
		}
		return &Station{}, stationErr
	}
	newStation := Station{}
	newStation.Name = station.Name
	newStation.Latitude = station.Latitude
	newStation.Longitude = station.Longitude
	newStation.Passengers = passengers
	res = collection.FindOneAndUpdate(context.Background(), temp, bson.M{
		"$set": bson.M{
			"name":       newStation.Name,
			"lat":        newStation.Latitude,
			"lng":        newStation.Longitude,
			"passengers": newStation.Passengers,
		}})

	resDecoded := Station{}
	err := res.Decode(&resDecoded)
	if err != nil {
		stationErr := &utils.AppErrors{
			Message:    "an error occured during decode the station to update v2",
			StatusCode: http.StatusBadGateway,
			Code:       "decode error",
		}
		return &Station{}, stationErr
	}
	return &newStation, nil
}

func GetStationByName(name string, client *mongo.Client) (*Station, *utils.AppErrors) {
	collection := client.Database("User_passenger").Collection("stations")
	temp := bson.M{"name": bson.M{"$eq": name}}
	res := collection.FindOne(context.Background(), temp)
	var station Station
	if err := res.Decode(&station); err != nil {
		stationErr := &utils.AppErrors{
			Message:    "an error occured during decode the station to update",
			StatusCode: http.StatusBadGateway,
			Code:       "decode error",
		}
		return &Station{}, stationErr
	}
	return &station, nil
}
