package domain

type User struct {
	Email    string  `bson:"email"`
	Name     string  `bson:"name"`
	Password string  `bson:"password"`
	Tokens   []Token `bson:"token"`
}

type Token struct {
	Context string `bson:"context"`
}
