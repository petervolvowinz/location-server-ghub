package main

import (
	"container/list"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/quaternion"
	"queue"
	"time"
	log "github.com/sirupsen/logrus"
	//"fmt"
	"io"
	"net/http"
)


var detectionQueue *list.List

func handler_c(w http.ResponseWriter, req *http.Request) {

	json := req.FormValue("gps")
	log.Println(json)

	if (queue.IsValidGPSJsonObject(json)){

		var riderObject = queue.GetGPSLocationObjectFromJSON(json)
		queue.AddNewPosition(riderObject)

		//go func(w  http.ResponseWriter) {
		warninglist := queue.RetrieveCollisionList(riderObject)
		ajsonlist := queue.RetrieveJSONList(warninglist)
		io.WriteString(w,ajsonlist)
		log.Info(" sent to client ", ajsonlist)
		//}(w)
	}

}

/*
{
"lat":<val>  latitude
"long":<val> longitude
"time":<val> timespan in seconds
"dist":<val> distance in meters
}
 */
func handler_w(w http.ResponseWriter, req *http.Request) {

	json := req.FormValue("search")
	log.Println(json)

	valid, searchObject := queue.IsValidSearchJsonObject(json)
	if (valid){
		fakeGPSSearchObject := &queue.GPSLocation{
			Location:  queue.Locationdata{
				Latitude:searchObject.Latitude,
				Longitude:searchObject.Longitude,
				Accuracy:1,
				Payload:"{}",
			},
			Gpsobject: 0,
			Uuid:      uuid.UUID{},
			Timestamp: 0,
		}

		log.Info(" searching for vehciles within: ",searchObject.Distance )
	}

}



func initObjectVehicleDetectionServer(){
	detectionQueue = queue.GetQueue()
}

// every 15 seconds dispose...
func Dispose(){
	for {
		queue.Out()
		var sleep = time.Second*queue.Expirationtime
		time.Sleep(sleep)
	}
}


func doAtimeTest(){

	time.Sleep(100 * time.Millisecond)
	ts1 := time.Now().UnixNano()
	log.Println(ts1)

	time.Sleep((1000 * time.Millisecond))
	ts2 := time.Now().UnixNano()
	log.Println(ts2)


	time.Sleep(5*time.Second)
	ts3 := time.Now().UnixNano()

	log.Println(ts3)

	log.Println( (ts2 - ts1) / 1e+9)
	log.Println( (ts3 - ts1) / 1e+9)

}

func main() {
    fmt.Println("starting server ...")

	initObjectVehicleDetectionServer()
	go Dispose()
	http.HandleFunc("/addposition", handler_c)
    http.HandleFunc("/retrieve",handler_w)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
