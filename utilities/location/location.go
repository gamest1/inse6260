package location

import (
	"fmt"
)

type Location struct {
  Latitude  float64 `bson:"lat" json:"lat"`
  Longitude float64 `bson:"long" json:"long"`
  Apartment string  `bson:"apartment" json:"apartment"`
  Number    int     `bson:"number" json:"number"`
  Street    string  `bson:"street" json:"street"`
  City      string  `bson:"city" json:"city"`
  State     string  `bson:"state" json:"state"`
  Zip       string  `bson:"zip" json:"zip"`
}

// DisplayLocation pretty prints Location.
func (l *Location) DisplayLocation() string {
	return fmt.Sprintf("Address: %d %s, %s %s  %s", l.Number,l.Street,l.City,l.State,l.Zip)
}

func (l *Location) String() string {
	return l.DisplayLocation()
}
