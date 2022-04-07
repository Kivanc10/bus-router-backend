package services

import (
	"bus-router-backend/mvc/domain"
	"bus-router-backend/mvc/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient = ConnectToMongoDb()

func ConnectToMongoDb() *mongo.Client {
	return domain.ConnectToMongoDb()
}

func SignUp(name, email, password string, client *mongo.Client) (*domain.User, *utils.AppErrors) {
	return domain.SignUp(name, email, password, client)
}

func ProcessJSONforUser(body []byte) (*domain.User, error) {
	return domain.ProcessJSONforUser(body)
}
func SignIn(name, password string, client *mongo.Client) (*domain.User, *utils.AppErrors) {
	return domain.SignIn(name, password, client)
}
func GetInside(client *mongo.Client) (*domain.User, *utils.AppErrors) {
	return domain.GetInside(client)
}
func LogoutForUser() *utils.AppErrors {
	return domain.LogoutForUser()
}
