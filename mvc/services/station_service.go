package services

import (
	"bus-router-backend/mvc/domain"
	"bus-router-backend/mvc/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

func AddStation(name string, lat, lng float64, passengers int64, client *mongo.Client) (*domain.Station, *utils.AppErrors) {
	return domain.AddStation(name, lat, lng, passengers, client)
}

func ProcessJSONforStation(body []byte) (*domain.Station, error) {
	return domain.ProcessJSONforStation(body)
}

func DeleteStation(name string, client *mongo.Client) (domain.Station, *utils.AppErrors) {
	return domain.DeleteStation(name, client)
}

func UpdatePassengerOfStation(name string, passengers int64, client *mongo.Client) (*domain.Station, *utils.AppErrors) {
	return domain.UpdatePassengerOfStation(name, passengers, client)
}

func GetStationByName(name string, client *mongo.Client) (*domain.Station, *utils.AppErrors) {
	return domain.GetStationByName(name, client)
}
