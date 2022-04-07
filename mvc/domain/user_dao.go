package domain

import (
	"bus-router-backend/mvc/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// mongodb bağla
// birbirine bağlı fonk yaz
// docker build -t yourusername/repository-name .
// docker run -d -p80:3000 yourusername/example-node-app
var mySigningKey = []byte("captainjacksparrowsayshi")

func ConnectToMongoDb() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017/")) //write 127.0.0.1 insted mongo for localhost
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil) // to check there is an error caused by interrupting
	if err != nil {
		panic(err)
	}
	return client
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isAlreadyExist(name string, client *mongo.Client) (User, bool) {
	collection := client.Database("User_passenger").Collection("user")
	temp := bson.M{"name": bson.M{"$eq": name}}
	result := User{}
	err := collection.FindOne(context.Background(), temp).Decode(&result)
	if err != nil {
		return User{}, false
	}
	return result, true
}

func createToken(name, email, password string) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", string(mySigningKey)) //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	//atClaims["user_id"] = userId
	atClaims["user_name"] = name
	atClaims["email"] = email
	atClaims["password"] = password
	atClaims["exp"] = time.Now().Add(time.Minute * 1500).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", errors.New("an error occured during the create token")
	}
	fmt.Println("jwt map --> ", atClaims)
	return token, nil
}

func addTokenToPerson(user *User) string {
	token, err := createToken(user.Name, user.Email, user.Password)
	if err != nil {
		log.Fatal("An error occured during the produce token ", err)
	}
	temp := Token{Context: token}
	user.Tokens = append(user.Tokens, temp)
	// keep it safe for login endpoint
	//os.Setenv("Token", user.Tokens[len(user.Tokens)-1].Context)
	//os.Setenv("username", user.Name)
	os.Setenv("last_registered_user", user.Name)
	return token
}

func insertInto(client *mongo.Client, user *User) *utils.AppErrors {
	collection := client.Database("User_passenger").Collection("user")
	hash, err := hashPassword(user.Password)
	if err != nil {
		log.Fatal("an error occured during hashed the password")
	}
	user.Password = hash
	//

	addTokenToPerson(user)
	toSave, err := bson.Marshal(user)
	if err != nil {
		log.Fatal("an error occured during the marshalling")
		return &utils.AppErrors{
			Message:    "an error occured during the marshalling",
			StatusCode: http.StatusNotFound,
			Code:       "marshall error",
		}
	}
	if _, result := isAlreadyExist(user.Name, client); !result {
		res, err := collection.InsertOne(context.Background(), toSave)
		if err != nil {
			panic(err)
		}
		id := res.InsertedID
		fmt.Println("id --> ", id)
		fmt.Println("user.tokens --> ", user.Tokens)
		return nil
	} else {
		return &utils.AppErrors{
			Message:    "the user is already exist",
			StatusCode: http.StatusUnauthorized,
			Code:       "already exist user is provided",
		}
	}
}
func SignUp(name, email, password string, client *mongo.Client) (*User, *utils.AppErrors) {
	var user User
	user.Name = name
	user.Email = email
	user.Password = password
	//addTokenToPerson(&user)
	if err := insertInto(client, &user); err != nil {
		log.Println("an error occured during the inserting ", err)
		return &User{}, err
	}
	return &user, nil
}

func ProcessJSONforUser(body []byte) (*User, error) {
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return &User{}, err
	}
	return &user, nil
}

// sign in

func addTokenForLogin(client *mongo.Client, user *User) *utils.AppErrors {
	collection := client.Database("User_passenger").Collection("user")
	filter := bson.M{"name": user.Name}
	update := bson.M{"$set": bson.M{
		"tokens": user.Tokens}}
	res := collection.FindOneAndUpdate(context.Background(), filter, update)
	resDecoded := User{}
	err := res.Decode(&resDecoded)
	if err != nil {
		return &utils.AppErrors{
			Message:    err.Error(),
			StatusCode: http.StatusNotAcceptable,
			Code:       "token update error",
		}
	}
	return nil

}

func SignIn(name, password string, client *mongo.Client) (*User, *utils.AppErrors) {
	collection := client.Database("User_passenger").Collection("user")
	filter := bson.M{"name": bson.M{"$eq": name}}
	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		userErr := &utils.AppErrors{
			Message:    "First sign up to sign in - there is no user exist has that name",
			StatusCode: http.StatusUnauthorized,
			Code:       "unauthorized",
		}
		return &User{}, userErr
	}
	if !checkPassword(password, user.Password) {
		userErr := *&utils.AppErrors{
			Message:    "Passwords are not matched,wrong password is provided",
			StatusCode: http.StatusUnauthorized,
			Code:       "wrong password",
		}
		return &User{}, &userErr
	}
	fmt.Println("point 0")
	addTokenToPerson(&user)
	fmt.Println("point 1")
	fmt.Println("user(p1) --> ", user)
	if err := addTokenForLogin(client, &user); err != nil {
		fmt.Println("point 2")
		return &User{}, err
	}

	os.Setenv("Token", user.Tokens[len(user.Tokens)-1].Context)
	os.Setenv("username", user.Name)
	fmt.Println("point 3")
	return &user, nil
}

func GetInside(client *mongo.Client) (*User, *utils.AppErrors) {
	if os.Getenv("username") == "" {
		userErr := *&utils.AppErrors{
			Message:    "please authenticate first",
			StatusCode: http.StatusUnauthorized,
			Code:       "unauthorized",
		}
		return &User{}, &userErr
	}
	user, err := findUserByName(os.Getenv("username"), client)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func findUserByName(name string, client *mongo.Client) (*User, *utils.AppErrors) {
	collection := client.Database("User_passenger").Collection("user")
	filter := bson.M{"name": bson.M{"$eq": name}}
	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		userErr := &utils.AppErrors{
			Message:    "First sign up to sign in - there is no user exist has that name",
			StatusCode: http.StatusUnauthorized,
			Code:       "unauthorized",
		}
		return &User{}, userErr
	}
	return &user, nil
}

func LogoutForUser() *utils.AppErrors {
	if os.Getenv("Token") != "" && os.Getenv("username") != "" {
		adminErr := &utils.AppErrors{
			Message:    "user is still logged in",
			StatusCode: http.StatusBadRequest,
			Code:       "unable to logout",
		}
		return adminErr
	}
	return nil
}
