package queue

import (
	"encoding/json"
	"github.com/google/uuid"
	"math"
	"testing"
)

func TestGetApproxDistance1(t *testing.T) {

	 lat1 := 37.387401
	 long1 := -122.035179

	 lat2 := 37.389649
	 long2 := -122.034433

	 expected_result := 260.0

	 outcome := GetApproxDistance1(lat1,long1,lat2,long2)

	if ( ( math.Abs(math.RoundToEven(outcome)) - expected_result) > 5){
	 	t.Error("Expected ",260,"+-5 got ",outcome)
	 }

}


func TestGetApproxDistance2(t *testing.T) {
	lat1 := 37.387401
	long1 := -122.035179

	lat2 := 37.389649
	long2 := -122.034433

	expected_result := 260.0

	outcome := GetApproxDistance2(lat1,long1,lat2,long2)

	if ( ( math.Abs(math.RoundToEven(outcome)) - expected_result) > 5){
		t.Error("Expected ",260,"+-5 got ",outcome)
	}
}


func TestGetRadians(t *testing.T) {
	pihalf := math.Pi / 2

	testval := GetRadians(90)

	if ( pihalf != testval ){
		t.Error("Expected radians(90) to be PI/2")
	}
}

func TestGetJSONFromGPSLocationObject(t *testing.T) {

	gps := &GPSLocation{
		Location: Locationdata{
			Latitude:37.387401,
			Longitude:-122.035179,
			Accuracy:1,
	    },
		Gpsobject:Bike,
		Uuid:uuid.New(),
		Timestamp:1,
	}

	ajson := GetGPSLocationJSON(*gps)
	bytes := []byte(string(ajson))

	gpsobject := &GPSLocation{}
	err := json.Unmarshal(bytes,gpsobject)

    if err != nil{
    	t.Error( "json could not unmarshal", err )
	}
	if (gps.Uuid != gpsobject.Uuid){
		t.Error(" json convertion failed ",  ajson)
	}

}

func TestUUID(t *testing.T) {

	uuid := GetUUID()
	gps := &GPSLocation{
		Location: Locationdata{
			Latitude:37.387401,
			Longitude:-122.035179,
			Accuracy:1,
		},
		Gpsobject:Bike,
		Uuid:uuid,
		Timestamp:1,
	}

	if (gps.Uuid != uuid){
		t.Error(" error generating uuid ")
	}

	uuid2 := GetUUID()
    if (uuid2 == gps.Uuid){
    	t.Error("uuid not unique")
	}
}



