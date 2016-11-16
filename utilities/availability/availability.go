package availability

import (
	"fmt"
)

type Availability struct {
  Monday    int `bson:"monday" json:"monday"`
  Tuesday   int `bson:"tuesday" json:"tuesday"`
  Wednesday int `bson:"wednesday" json:"wednesday"`
  Thursday  int `bson:"thursday" json:"thursday"`
  Friday    int `bson:"friday" json:"friday"`
  Saturday  int `bson:"saturday" json:"saturday"`
  Sunday    int `bson:"sunday" json:"sunday"`
}

// DisplayAvailability displays availability in human-readable format.
func (availability *Availability) DisplayAvailability() string {
	return fmt.Sprintf("Hours:\n\nMonday %d\nTuesday %d\nWednesday %d...", availability.Monday, availability.Tuesday, availability.Wednesday)
}

func (a *Availability) String() string {
	return a.DisplayAvailability()
}
