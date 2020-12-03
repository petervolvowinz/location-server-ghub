package queue

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type featureCollection struct {
	Type     string    `json:"type"`
	Features []feature `json:"features"`
}

type feature struct {
	Type       string     `json:"type"`
	Properties properties `json:"properties"`
	Geometry   geometry   `json:"geometry"`
}

type properties struct {
	Ambtemp    float64   `json:"ambienttemp"`
	Cabintemp  float64   `json:"cabintemp"`
	Drivertemp float64   `json:"drivertemp"`
	Day        string    `json:"Day"`
	Time       string    `json:"Time"`
	Icontype   string    `json:"Icontype"`
	UUID       uuid.UUID `json:"UUID"`
}

type geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func getDayStr(ti time.Time) string {

	switch ti.Weekday() {
	case time.Friday:
		return "Fri"
	case time.Saturday:
		return "Sat"
	case time.Sunday:
		return "Sun"
	case time.Monday:
		return "Mon"
	case time.Tuesday:
		return "Tue"
	case time.Wednesday:
		return "Wed"
	case time.Thursday:
		return "Thu"
	default:
		return "NIL"
	}
}

func getTimeStr(ti time.Time) string {
	return ti.UTC().String()
}

func getIconStr(vehicletype int) string {
	switch vehicletype {
	case Bike:
		return "bicycle"
	case Car:
		return "car"
	case Pedestrian:
		return "pedestrian"
	case Truck:
		return "truck"
	case Skater:
		return "skater"
	default:
		return "raccoon"
	}
}

func getDayIconTypeTime(gps GPSLocation) (string, string, string) {
	tm := time.Unix(gps.Timestamp, 0)

	dayS := getDayStr(tm)
	timeS := getTimeStr(tm)
	iconS := getIconStr(gps.Gpsobject)

	return timeS, dayS, iconS
}

func getTemp(gps GPSLocation) (float64, float64, float64) {
	switch gps.Gpsobject {
	case Car:
		var bytes = []byte(string(gps.Location.Payload))
		cpl := &Climatepayload{}
		err := json.Unmarshal(bytes, &cpl)
		if err != nil {
			log.Info("could not unmarshal payload ", err)
			return 0, 0, 0
		}
		return cpl.Ambientemp, cpl.Cabintemp, cpl.Drivertemp
	default:
		return 0, 0, 0
	}
}

func ConvertToMapBoxFreindlyJSON(hits []interface{}) string {

	var features []feature

	for _, element := range hits {
		gps := element.(GPSLocation)
		dayS, TimeS, IcontypeS := getDayIconTypeTime(gps)
		at, ct, dt := getTemp(gps)
		afeature := &feature{
			Type: "Feature",
			Properties: properties{
				Ambtemp:    at,
				Cabintemp:  ct,
				Drivertemp: dt,
				Day:        dayS,
				Time:       TimeS,
				Icontype:   IcontypeS,
				UUID:       gps.UI,
			},
			Geometry: geometry{
				Type:        "Point",
				Coordinates: []float64{gps.Location.Longitude, gps.Location.Latitude},
			},
		}
		features = append(features, *afeature)
	}

	bytes, err := json.Marshal(features)

	if err != nil {
		log.Info("could not marshal mapbox json  ")
		return ""
	}

	return string(bytes)
}
