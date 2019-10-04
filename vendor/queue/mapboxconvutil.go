package queue

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type featureCollection struct{
	Type string	`json:"type"`
	Features [] feature `json:"features"`
}

type feature struct{
	Type string `json:"type"`
	Properties properties `json:"properties"`
	Geometry geometry `json:"geometry"`
}

type properties struct{
	Ambtemp int `json:"ambienttemp"`
	Cabintemp int `json:"cabintemp"`
	Drivertemp int `json:"drivertemp"`
	Day string `json:"Day"`
	Time string `json:"Time"`
	Icontype string `json:"Icontype"`
	UUID uuid.UUID `json:"UUID"`
}

type geometry struct{
	Type string `json:"type"`
	Coordinates [] float64 `json:"coordinates"`
}


func ConvertToMapBoxFreindlyJSON(hits [] interface{}) string{

	var features [] feature
	
	for _,element := range hits {
		gps := element.(GPSLocation)
		afeature := &feature{
			Type: "Feature",
			 Properties: properties{
				Ambtemp:    70,
				Cabintemp:  75,
				Drivertemp: 75,
				Day:        "10",
				Time:       "10:00",
				Icontype:   "Car",
				UUID:       gps.Uuid,
			},
			Geometry: geometry{
				Type:        "geometry",
				Coordinates: []float64{gps.Location.Latitude,gps.Location.Longitude},
			},
		}
		features = append(features, *afeature)
	}

	bytes,err := json.Marshal(features)

	if (err != nil){
		log.Info("could not marshal mapbox json  ")
		return "";
	}

	return string(bytes)
}



