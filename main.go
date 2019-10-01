package main

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"queue"
	"time"
	//"fmt"
	"io"
	"net/http"
)



func handler_c(w http.ResponseWriter, req *http.Request) {

	json := req.FormValue("gps")
	//log.Println(json)

	if (queue.IsValidGPSJsonObject(json)){

		var riderObject = queue.GetGPSLocationObjectFromJSON(json)
		riderObject.Timestamp = time.Now().UnixNano() // timestamp as soon as we can.
		queue.AddNewPosition(riderObject)

		warninglist := queue.RetrieveCollisionList(riderObject)
		ajsonlist := queue.RetrieveJSONList(warninglist)
		io.WriteString(w,ajsonlist)

		//log.Info(" sent to client ", ajsonlist)
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

		log.Info(" searching for vehciles within: ",fakeGPSSearchObject.Location.Longitude )
	}

}

// to be able to check if it is alive from
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("☄ HTTP status code returned!"))
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

	go Dispose()

	http.HandleFunc("/addposition", handler_c)
    http.HandleFunc("/retrieve",handler_w)
    http.HandleFunc("/version",pingHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
