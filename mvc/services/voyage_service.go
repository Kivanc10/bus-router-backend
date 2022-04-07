package services

import (
	"bus-router-backend/mvc/domain"
	"bus-router-backend/mvc/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

func ProcessJSONforVoyage(body []byte) (*domain.Voyage, error) {
	return domain.ProcessJSONforVoyage(body)
}

func AddVoyage(from, to string, distance int64, client *mongo.Client) (*domain.Voyage, *utils.AppErrors) {
	return domain.AddVoyage(from, to, distance, client)
}

func DeleteVoyage(from, to string, client *mongo.Client) (domain.Voyage, *utils.AppErrors) {
	return domain.DeleteVoyage(from, to, client)
}
