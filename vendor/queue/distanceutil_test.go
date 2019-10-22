package queue

import (
	"encoding/json"
	"github.com/google/uuid"
	"math"
	"testing"
	"time"
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
		UI:uuid.New(),
		Timestamp:1,
	}

	ajson := GetGPSLocationJSON(*gps)
	bytes := []byte(string(ajson))

	gpsobject := &GPSLocation{}
	err := json.Unmarshal(bytes,gpsobject)

    if err != nil{
    	t.Error( "json could not unmarshal", err )
	}
	if (gps.UI != gpsobject.UI){
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
		UI :uuid,
		Timestamp:1,
	}

	if (gps.UI != uuid){
		t.Error(" error generating uuid ")
	}

	uuid2 := GetUUID()
    if (uuid2 == gps.UI){
    	t.Error("uuid not unique")
	}
}


/*
{
    "type": "FeatureCollection",
    "features": [{
        "type": "Feature",
        "properties": {
            "ambienttemp": 80,
            "cabintemp": 75,
            "drivertemp": 68,
            "Day": "Fri",
            "Time": "10:00",
            "Icontype": "Car",
            "UUID": "dfdfdf"
        },
        "geometry": {
            "type": "Point",
            "coordinates": [-122.0349794626236, 37.387971267871]
        }
    }, {
        "type": "Feature",
        "properties": {
            "ambienttemp": 80,
            "cabintemp": 75,
            "drivertemp": 68,
            "Day": "Fri",
            "Time": "10:00",
            "Icontype": "Car",
            "UUID": "dfdfdf"
        },
        "geometry": {
            "type": "Point",
            "coordinates": [-122.0349794626236, 37.387971267871]
        }
    }]
}
 */

func gettestWebJSON() string {
	return `{
		"type": "FeatureCollection",
			"features": [{
			"type": "Feature",
			"properties": {
				"ambienttemp": 80,
				"cabintemp": 75,
				"drivertemp": 68,
				"Day": "Fri",
				"Time": "10:00",
				"Icontype": "Car",
				"UUID": "2d998dc6-0b66-4d27-aeb5-dccbd73489c1"
			},
			"geometry": {
				"type": "Point",
				"coordinates": [-122.0349794626236, 37.387971267871]
			}
		}, {
			"type": "Feature",
				"properties": {
				"ambienttemp": 80,
					"cabintemp": 75,
					"drivertemp": 68,
					"Day": "Fri",
					"Time": "10:00",
					"Icontype": "Car",
					"UUID": "2d998dc6-0b66-4d27-aeb5-dccbd73489c1"
			},
			"geometry": {
				"type": "Point",
					"coordinates": [-122.0349794626236, 37.387971267871]
}
}]
}`

}
func TestWebJSON(t *testing.T){

	features := &feature{
		Type: "Feature",
		Properties: properties{
			Ambtemp:    80,
			Cabintemp:  75,
			Drivertemp: 68,
			Day:        "Fri",
			Time:       "10:00",
			Icontype:   "Car",
			UUID:       uuid.New(),
		},
		Geometry: geometry{
			Type:"Point",
			Coordinates:[]float64{-122.0349794626236,
				37.387971267871},
		},
	}

	featurecollection := &featureCollection{
		Type: "FeatureCollection",
		Features:          []feature{*features,*features},
	}


	bytes,err := json.Marshal(featurecollection)
	if (err != nil){
		t.Error(" json not marshalled ", err)
	}

	jsonstring := gettestWebJSON()
	bytes = []byte(string(jsonstring))

	featurecollectionfromjsonstring := &featureCollection{}
	err = json.Unmarshal(bytes,featurecollectionfromjsonstring)

	if (err != nil){
		t.Error(" json string not unmarshalled to featureCollection ", err)
	}
	if (featurecollectionfromjsonstring.Type != featurecollection.Type){
		t.Error(" featurecollection josn to object is not marshalled correctly " )
	}


}

func getwarningsJSON() string {
	return `{"warnings":[{"Location":{"lat":37.390750000000104,"lng":-122.03407102774996,"accuracy":1,"payload":"{\"ambientemp\":23.3,\"cabintemp\":19.7,\"drivertemp\":22,\"parkingspots\":0}"},"Gpsobject":1,"UUID":"78401d36-716d-4670-9a90-18b0483f94e4","Timestamp":1570752510331861000}]}`
}

func TestResultList(t *testing.T){
	jsonStr := getwarningsJSON()

	obj := &Warninglst{}
	err := json.Unmarshal([]byte(jsonStr),obj)

	if (err != nil){
		t.Error(" json string not unmarshalled to Warninglst ", err)
	}
}


func TestTimeSpanFunc(t *testing.T){

	ts1 := time.Now().UnixNano()
	time.Sleep(1*time.Second)
	ts2 := time.Now().UnixNano()

	if withinTimeSpan(ts1,ts2,Timedepth) == false{
		t.Error(" 1 second should be within timespan")
	}
	time.Sleep(5*time.Second)
	ts2 = time.Now().UnixNano()
	if withinTimeSpan(ts1,ts2,Timedepth){
		t.Error(" 6 seconds should be outside timespan")
	}

}


func TestWithinDistanceFunc(t *testing.T){
	lat1 := 37.387401
	long1 := -122.035179

	loc1 := GPSLocation{Location:Locationdata{Latitude:lat1,Longitude:long1,Accuracy:0,Payload:"{}"}}

	lat2 := 37.389649
	long2 := -122.034433

	loc2 := GPSLocation{Location:Locationdata{Latitude:lat2,Longitude:long2,Accuracy:0,Payload:"{}"}}

	if withinDistance(loc1,loc2,300) == false{
		t.Error("expected within distance")
	}

	if withinDistance(loc2,loc1,300) == false{
		t.Error("expected within distance")
	}

	if withinDistance(loc1,loc2,250){
		t.Error("expected without distance")
	}

	if withinDistance(loc2,loc1,250){
		t.Error("expected without distance")
	}
}