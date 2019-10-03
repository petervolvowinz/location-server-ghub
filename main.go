package main

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"queue"
	"time"
	//"fmt"
	"io"
	"net/http"
)



func handleNewEntries(w http.ResponseWriter, req *http.Request) {

	log.Info("trying to connect")

	json := req.FormValue("gps")

	valid,riderObject := queue.IsValidGPSLocationJSON(json)
	log.Info(json)
	if (valid){

		riderObject.Timestamp = time.Now().UnixNano() // timestamp as soon as we can.
		queue.AddNewPosition(*riderObject)

		warninglist := queue.RetrieveCollisionList(*riderObject)
		ajsonlist := queue.GetWarninglistJSON(warninglist)

        w.WriteHeader(http.StatusOK)
		io.WriteString(w,ajsonlist)

	}else {
		log.Info("Invalid json sent from client")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w,"server could not parse json parameter")
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
func handleGPSFence(w http.ResponseWriter, req *http.Request) {

	json := req.FormValue("search")
	log.Println(json)

	valid, searchObject := queue.IsValidSearchstructJSON(json)
	if (valid){
		GPSSearchObject := &queue.GPSLocation{
			Location:  queue.Locationdata{
				Latitude:searchObject.Latitude,
				Longitude:searchObject.Longitude,
				Accuracy:1,
				Payload:"{}",
			},
			Gpsobject: 0,
			Uuid:      uuid.UUID{},
			Timestamp: time.Now().UnixNano(),
		}

		//timespan := searchObject.Timespan
		//distance := searchObject.Distance


		log.Info(" searching for vehciles within: ",GPSSearchObject.Location.Longitude )

	}

}

// to be able to check if it is alive from
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(" HTTP status code returned!"))
}


// every 100 ms seconds dispose...
func Dispose(){
	for {
		queue.Out()
		var sleep = time.Millisecond*100
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
    log.Info("starting server ...")

	go Dispose()

	http.HandleFunc("/addposition", handleNewEntries)
    http.HandleFunc("/retrieve",handleGPSFence)
    http.HandleFunc("/version",pingHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))

    log.Info("closing server down ... ")
}
