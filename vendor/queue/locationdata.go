package queue

import (
	"container/list"
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var queueMutex = &sync.Mutex{}
var once sync.Once

const (

	MaxDetections = 5
	Expirationtime = 5 // throw away any entries older that this (seconds)
	Timedepth = 5 // only consider queue neighbours with in this (seconds)

	Criticaldistance = 200 // The distance to a bike/car where we issue a warning (meters)
	GARBAGESIZE = 5
)

const (
	Bike int = iota
	Car
	Raccoon
)

type Warninglst struct{
	Warnings []GPSLocation	`json:warnings`
}

type GPSLocation struct{
	Location Locationdata `json:location`
	Gpsobject int	  `json:gpsobject`
	Uuid uuid.UUID       `json:uuid`
	Timestamp int64    `json:timestamp`
}

type Locationdata struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Accuracy  float64 `json:"accuracy"`
	Payload   string  `json:"payload"`
}

type test struct{
	More []Locationdata `json:more`
}

type Climatepayload struct{
	Ambientemp float64 `json:"ambientemp"`
	Cabintemp float64 `json:"cabintemp"`
	Drivertemp float64 `json:"drivertemp"`
}

type SearchJSON struct{
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Distance  int `json:"dist"`
	TimeSpan  int `json:"timespan"`
}

func GetClimateJSON(climate Climatepayload) string{
	bytes,err := json.Marshal(climate)
	if (err != nil){
		log.Info("could not marshal climate object into json")
		return ""
	}
	return string(bytes)
}

var(
	instance *list.List
)

var(
	filter *Filtervalues
)

// singleton
func GetQueue() *list.List{

	once.Do(func(){
		instance = list.New()
		filter = SetFilterValue(GetFilterValues(),Timedepth,Criticaldistance)
	})

	return instance
}

func AddNewPosition(location GPSLocation){
	location.Timestamp = time.Now().UnixNano()
	//protect queue
	queueMutex.Lock()
	GetQueue().PushFront(location)

	queueMutex.Unlock()
}

// helper functions to collect object warning list, e.g bikes or vehicles.
func withinTime(driver_ts int64,detect_ts int64) bool{
	return ((driver_ts - detect_ts)/1e+9 < int64(filter.timespanvalue))
}

func withinDistance(driver GPSLocation,detect GPSLocation) bool{
	lat1 := driver.Location.Latitude
	long1 := driver.Location.Longitude

	lat2 := detect.Location.Latitude
	long2 := detect.Location.Longitude

	dist := GetApproxDistance2(lat1,long1,lat2,long2)

	return (dist < float64(filter.distancevalue))
}

/*func nearbyObject2(driver GPSLocation,detect GPSLocation,vehicletype int,distCheck,timeCheck) bool){

}*/

func nearbyObject(driver GPSLocation,detect GPSLocation,vehicletype int) bool{
	// if it is the same don't add
	if (driver.Uuid == detect.Uuid) {
		return false
	}

	// don't need to check type its going to be mutual exclusive
	if (withinTime(driver.Timestamp,detect.Timestamp)){
		if withinDistance(driver,detect){
			return true
		}
	}


	return false
}

func SetFilter(fv Filtervalues){
	filter = &fv
}

func RetrieveCollisionList(objecttype GPSLocation)[]GPSLocation{

	queueMutex.Lock()

	var listofdectees []GPSLocation

    list := GetQueue()

	for element := list.Front();element != nil;element = element.Next(){
		var detectee = element.Value.(GPSLocation)
		if nearbyObject(objecttype,detectee,objecttype.Gpsobject){
			listofdectees = append(listofdectees, detectee)
			log.Info(" collisions " , len(listofdectees))
		}

	}

	queueMutex.Unlock()

	return listofdectees
}

//Garbage collection
func Out(){
	queueMutex.Lock()

	Q := GetQueue()
	for i := 0; i < GARBAGESIZE; i++{ // DO GARBAGESIZe removals.
		if Q.Len() > 0 {
			element := Q.Back()

			item := element.Value.(GPSLocation)
			if retieree(item) {
				log.Info("expired : ", item.Timestamp)
				Q.Remove(element)
			}

		} else {
			break;
		}
	}
	log.Println("Q SIZE IS : ", Q.Len())
	queueMutex.Unlock()
}


func retieree(oldie GPSLocation) bool{

	currentsecond := time.Now().UnixNano()
	oldiesecond := oldie.Timestamp

	log.Println("time checking ",currentsecond - oldiesecond)
	return (currentsecond -  oldiesecond)/1e+9 > Expirationtime

}

func GetJSONFromGPSLocationObject(obj GPSLocation) string{

	result,err := json.Marshal(obj)
	if (err != nil){
		log.Error(" json not marshalled ", err)
	}

	return string(result)
}

// getUUID
func getUUID()(uuid.UUID){
	 return uuid.New()
}
//GetGPSLocationObjectFromJSON
func GetGPSLocationObjectFromJSON(ajson string) GPSLocation{

	var bytes = []byte(string(ajson))
	gps := &GPSLocation{}
	err := json.Unmarshal(bytes,&gps)
	if (err != nil){
		log.Error("invalid json ", err)
	}

    gps.Uuid = getUUID()
	return *gps
}


func IsValidGPSJsonObject(ajson string) bool{

	var bytes = []byte(ajson)
	gps := &GPSLocation{}
	err := json.Unmarshal(bytes,&gps)

	return (err == nil)
}

func RetrieveJSONList(warninglist []GPSLocation) string{

	warnings := &Warninglst{warninglist}
	result,err := json.Marshal(warnings)


	if (err != nil){
		log.Error("could not build json GPSLocation list", err)
		return "{error:" + "json list generation failed}"
	}

	return string(result)
}
