package services

import (
	"bus-router-backend/mvc/domain"
	"bus-router-backend/mvc/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAdmin(adminCode int64) (*domain.Admin, *utils.AppErrors) {
	return domain.GetAdmin(adminCode)
}

func AdminLogin(admin *domain.Admin) (*domain.Admin, *utils.AppErrors) {
	return domain.AdminLogin(admin)
}

func ProcessJSONforAdmin(body []byte) (*domain.Admin, error) {
	return domain.ProcessJSONforAdmin(body)
}

func AdminInside(admin domain.Admin) *utils.AppErrors {
	return domain.AdminInside(admin)
}

func AdminLogout(admin domain.Admin) *utils.AppErrors {
	return domain.AdminLogout(admin)
}

func GetAllUsers(client *mongo.Client) ([]bson.M, *utils.AppErrors) {
	return domain.GetAllUsers(client)
}

func GetAllStations(client *mongo.Client) ([]bson.M, *utils.AppErrors) {
	return domain.GetAllStations(client)
}

func GetAllVoyages(client *mongo.Client) ([]bson.M, *utils.AppErrors) {
	return domain.GetAllVoyages(client)
}
