package domain

type Admin struct {
	Name     string `bson:"name"`
	Password string `bson:"password"`
}
