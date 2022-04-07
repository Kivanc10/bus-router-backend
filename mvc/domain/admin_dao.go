package domain

import (
	"bus-router-backend/mvc/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	admins = map[int64]*Admin{
		123: &Admin{
			Name:     "root",
			Password: "root",
		},
	}
)

func GetAdmin(adminCode int64) (*Admin, *utils.AppErrors) {
	if admin := admins[adminCode]; admin != nil {
		return admin, nil
	} else {
		return nil, &utils.AppErrors{
			Message:    "there is no admin that have that code",
			StatusCode: http.StatusNotFound,
			Code:       "not_found",
		}
	}
}

func AdminLogin(admin *Admin) (*Admin, *utils.AppErrors) {
	if admin.Name != "root" || admin.Password != "root" {
		adminErr := &utils.AppErrors{
			Message:    "unable to login as admin,provided username or root is wrong",
			StatusCode: http.StatusNotFound,
			Code:       "no content",
		}
		fmt.Println(adminErr.Message)
		return &Admin{}, adminErr
	}
	os.Setenv("admin", admin.Name)
	return admin, nil
}

func ProcessJSONforAdmin(body []byte) (*Admin, error) {
	var admin Admin
	if err := json.Unmarshal(body, &admin); err != nil {
		return &Admin{}, err
	}
	return &admin, nil
}

func AdminInside(admin Admin) *utils.AppErrors {
	//middleware.MiddlewareForAdmin(http.HandlerFunc(next))
	if admin.Name != "root" || admin.Password != "root" {
		adminErr := *&utils.AppErrors{
			Message:    "username or password is wrong",
			StatusCode: http.StatusNotFound,
			Code:       "wrong credentials are provided",
		}
		return &adminErr
	}
	return nil
}

func AdminLogout(admin Admin) *utils.AppErrors {
	if os.Getenv("admin") != "" {
		adminErr := *&utils.AppErrors{
			Message:    "Admin is still logged in",
			StatusCode: http.StatusNotFound,
			Code:       "unable to log out",
		}
		return &adminErr
	}
	return nil
}

func GetAllUsers(client *mongo.Client) ([]bson.M, *utils.AppErrors) {
	collection := client.Database("User_passenger").Collection("user")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return []bson.M{}, &utils.AppErrors{
			Message:    "an error occured during fetch the all results belong to the users",
			StatusCode: http.StatusBadGateway,
			Code:       "mongo fetch error",
		}
	}
	var users []bson.M
	if err = cursor.All(context.Background(), &users); err != nil {
		adminErr := &utils.AppErrors{
			Message:    "an error occured during fetch result from cursor",
			StatusCode: http.StatusBadRequest,
			Code:       "cursor error",
		}
		return []bson.M{}, adminErr
	}
	return users, nil
}

func GetAllStations(client *mongo.Client) ([]bson.M, *utils.AppErrors) {
	collection := client.Database("User_passenger").Collection("stations")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return []bson.M{}, &utils.AppErrors{
			Message:    "an error occured during fetch the all results belong to the stations",
			StatusCode: http.StatusBadGateway,
			Code:       "mongo fetch error",
		}
	}
	var stations []bson.M
	if err = cursor.All(context.Background(), &stations); err != nil {
		adminErr := &utils.AppErrors{
			Message:    "an error occured during fetch result from cursor",
			StatusCode: http.StatusBadRequest,
			Code:       "cursor error",
		}
		return []bson.M{}, adminErr
	}
	return stations, nil
}

func GetAllVoyages(client *mongo.Client) ([]bson.M, *utils.AppErrors) {
	collection := client.Database("User_passenger").Collection("voyages")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return []bson.M{}, &utils.AppErrors{
			Message:    "an error occured during fetch the all results belong to the voyages",
			StatusCode: http.StatusBadGateway,
			Code:       "mongo fetch error",
		}
	}
	var voyages []bson.M
	if err = cursor.All(context.Background(), &voyages); err != nil {
		adminErr := &utils.AppErrors{
			Message:    "an error occured during fetch result from cursor",
			StatusCode: http.StatusBadRequest,
			Code:       "cursor error",
		}
		return []bson.M{}, adminErr
	}
	return voyages, nil

}
