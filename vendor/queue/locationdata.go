package queue

import (
	"encoding/json"
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

var queueMutex = &sync.Mutex{}
var once sync.Once

const (
	Expirationtime   = 15  // throw away any entries older that this (seconds)
	Timedepth        = 5   // only consider queue neighbours with in this (seconds)
	Criticaldistance = 200 // The distance to a bike/car where we issue a warning (meters)
)

const (
	Bike int = iota
	Car
	Raccoon
)

var (
	instance *dll.List
	filter   *Filtervalues
)

type Searchstruct struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Distance  int64   `json:"distance"`
	Timespan  int64   `json:"timespan"`
}

type Warninglst struct {
	Warnings []interface{} `json:"warnings"`
}

type GPSLocation struct {
	Location  Locationdata `json:"Location"`
	Gpsobject int          `json:"Gpsobject"`
	UI        uuid.UUID    `json:"UUID"`
	Timestamp int64        `json:"timestamp"`
}

type Locationdata struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Accuracy  float64 `json:"accuracy"`
	Payload   string  `json:"payload"`
}

type Climatepayload struct {
	Ambientemp   float64 `json:"ambientemp"`
	Cabintemp    float64 `json:"cabintemp"`
	Drivertemp   float64 `json:"drivertemp"`
	Parkingspots int64   `json:"parkingspots"`
}

/*type SearchJSON struct{
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Distance  int `json:"dist"`
	TimeSpan  int `json:"timespan"`
}*/


func AddNewPosition(location GPSLocation) {
	addNewPosition_2(location)
}

func addNewPosition_2(location GPSLocation) {
	location.Timestamp = time.Now().UnixNano()

	queueMutex.Lock()

	Add(location) // add entry first...

	queueMutex.Unlock()
}


//obsolete as should be deleted
func RetrieveCollisionList(objecttype GPSLocation) []interface{} {

	return RetrieveCollisionList_2(objecttype, Timedepth, Criticaldistance)
}

func GetDefaultParams() (int64, int64) {
	return Timedepth, Criticaldistance
}

func RetrieveCollisionList_2(objecttype GPSLocation, timed int64, dist int64, depth ...int) []interface{} {

	queueMutex.Lock()

	timedistfilter := &TimeDistFilter{
		distance: dist,
		time:     timed,
	}

	var listofdectees []interface{}

	if len(depth) > 0 {
		listofdectees = FindAll(objecttype, timedistfilter, withinTimeAndDistanceFilter, depth[0])
	} else {
		listofdectees = FindAll(objecttype, timedistfilter, withinTimeAndDistanceFilter)
	}

	queueMutex.Unlock()

	return listofdectees
}

//Garbage collection, just pick the last one if it is eligeable
func Remove() {
	// RemoveOldData()
	removeLast()
	log.Info("Q SIZE IS : ", GetQueueInstance().Size())
}

func RemoveAllOld() {
	queueMutex.Lock()

	//Q := GetQueueInstance()
	index := Find()
	RemoveAll(index)

	queueMutex.Unlock()
}

func removeLast() {
	queueMutex.Lock()

	Q := GetQueueInstance()
	it := Q.Iterator()

	if it.Last() && retieree(it.Value().(GPSLocation)) {
		Q.Remove(Q.Size() - 1)
	}

	queueMutex.Unlock()
}

func retieree(oldie GPSLocation) bool {

	currentsecond := time.Now().UnixNano()
	oldiesecond := oldie.Timestamp

	log.Println("time checking ", (currentsecond-oldiesecond)/1e+9)
	return (currentsecond-oldiesecond)/1e+9 > Expirationtime

}

/***************************
 	JSON marshal and unmarshal
**************************/

// Returns a JSON string of a Climatepayload struct
func GetClimatepayloadJSON(climate Climatepayload) string {
	bytes, err := json.Marshal(climate)
	if err != nil {
		log.Info("could not marshal climate object into json")
		return ""
	}
	return string(bytes)
}

// Returns a JSON string of a GPSLocation JSON
func GetGPSLocationJSON(gps GPSLocation) string {

	bytes, err := json.Marshal(gps)
	if err != nil {
		log.Error(" json not marshalled ", err)
	}

	return string(bytes)
}

func GetGeneralJSON(str interface{}) string {

	bytes, err := json.Marshal(str)
	if err != nil {
		log.Error(" json not marshalled ", err)
	}

	return string(bytes)
}

//Returns a GPSLocation struct from JSON
func GetGPSLocationFromJSON(ajson string) GPSLocation {

	var bytes = []byte(string(ajson))
	gps := &GPSLocation{}
	err := json.Unmarshal(bytes, &gps)
	if err != nil {
		log.Error("invalid json ", err)
	}

	return *gps
}

// Returns true if JSON is a GPSLocation struct and returns the struct
func IsValidGPSLocationJSON(ajson string) (bool, *GPSLocation) {

	var bytes = []byte(ajson)
	gps := &GPSLocation{}
	err := json.Unmarshal(bytes, &gps)

	return (err == nil), gps
}

// Returns true if JSON is a Searchstruct struct and returns the struct
func IsValidSearchstructJSON(ajson string) (bool, *Searchstruct) {
	var bytes = []byte(ajson)
	srt := &Searchstruct{}
	err := json.Unmarshal(bytes, &srt)

	return (err == nil), srt
}

// Returns a json of a warningslist array of Warninglst structs
func GetWarninglistJSON(warninglist []interface{}) string {

	warnings := &Warninglst{warninglist}
	result, err := json.Marshal(warnings)

	if err != nil {
		log.Error("could not build json GPSLocation list ", err)
		return "{error:" + "json list generation failed}"
	}

	return string(result)
}

// Returns a new UUID
func GetUUID() uuid.UUID {
	return uuid.New()
}

// check params and convert from string
func ConvertTimeDistanceParams(time string, distance string) (int64, int64, error) {
	times, err := strconv.ParseInt(time, 10, 64)
	if err != nil {
		log.Info("time param not parsed")
		return -1, -1, err
	}
	dist, err := strconv.ParseInt(distance, 10, 64)
	if err != nil {
		log.Info("time param not parsed")
		return -1, -1, err
	}
	return times, dist, nil
}
