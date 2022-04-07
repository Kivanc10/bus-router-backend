package domain

type Station struct {
	Name       string  `bson:"name"`
	Latitude   float64 `bson:"lat"`
	Longitude  float64 `bson:"lng"`
	Passengers int64   `bson:"passengers"`
}
