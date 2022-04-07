package domain

type Voyage struct {
	From     Station `bson:"from"`
	To       Station `bson:"to"`
	Distance int64   `bson:"distance"`
}
