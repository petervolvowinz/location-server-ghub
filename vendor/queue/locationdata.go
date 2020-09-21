package queue

// do this
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
	Zindex    uint64  `json:zindex`
}

type Climatepayload struct {
	Ambientemp   float64 `json:"ambientemp"`
	Cabintemp    float64 `json:"cabintemp"`
	Drivertemp   float64 `json:"drivertemp"`
	Parkingspots int64   `json:"parkingspots"`
}

type LocationData interface{
	AddNewPosition(location GPSLocation)
	GetDefaultParams()(int64,int64)
	RetrieveCollisionList(objecttype GPSLocation, timed int64, dist int64, depth ...int) []interface{}
	Remove()
	Init()
}

type LocationDataHandler struct{
	users RoadUsers
}


var (
	locationdata RoadUsers
	locationdatainstantiator = RoadUsersFactory()
	datahandler LocationDataHandler
)

func GetLocationDataInstance() RoadUsers {
	locationdata = locationdatainstantiator("T")
	return locationdata
}

/****** Implementation of LocationData interface *********/

func (Lh *LocationDataHandler) Init(){
	Lh.users = GetLocationDataInstance()
}

func (Lh *LocationDataHandler) AddNewPosition(location GPSLocation) {
	queueMutex.Lock()

	location.Timestamp = time.Now().UnixNano()
	Lh.users.AddRoadUserPosition(location)

	queueMutex.Unlock()
}


func  (Lh *LocationDataHandler) GetDefaultParams() (int64, int64) {
	return Timedepth, Criticaldistance
}

func (Lh *LocationDataHandler) RetrieveCollisionList(objecttype GPSLocation, timed int64, dist int64, depth ...int) []interface{} {

	queueMutex.Lock()

	timedistfilter := &TimeDistFilter{
		distance: dist,
		time:     timed,
	}

	var listofdectees []interface{}

	breaklimit := 10
	if len(depth) > 0 {
		breaklimit = depth[0]
	}
	listofdectees = Lh.users.GetNearbyRoadUsers(objecttype,timedistfilter,withinTimeAndDistanceFilter,breaklimit)

	queueMutex.Unlock()

	return listofdectees
}

//Garbage collection, just pick the last one if it is eligeable
func   (Lh *LocationDataHandler) Remove() {
	queueMutex.Lock()
	Lh.users.GarbageCollect()
	queueMutex.Unlock()
}

/*func RemoveAllOld() {
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
}*/


/***************************
 	JSON marshal and unmarshal uti
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
