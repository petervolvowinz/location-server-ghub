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


func addPosition(json string){
	log.Info("trying to connect")

	valid,riderObject := queue.IsValidGPSLocationJSON(json)
	log.Info(json)
	if (valid) {
		riderObject.Timestamp = time.Now().UnixNano() // timestamp as soon as we can.
		queue.AddNewPosition(*riderObject)
	}
}

// just add , no response data
func addNoJsonResponse(w http.ResponseWriter, req *http.Request){
	json := req.FormValue("gps")
	go addPosition(json)
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("position inserted"))
}

//f08db884-f205-4153-9211-8b29245bbd89
func handleNewEntries(w http.ResponseWriter, req *http.Request) {

	log.Info("trying to connect")

	json := req.FormValue("gps")

	t_param,d_param,err := queue.ConvertTimeDistanceParams(req.FormValue("timespan"),req.FormValue("distance"))

	//log.Info("timespan is ", t_param , " distance span is ", d_param)

	if err != nil {
		log.Info("Invalid time and/or distance params  sent from client, using defaults as fallback")
		t_param,d_param = queue.GetDefaultParams()
	}

	valid,riderObject := queue.IsValidGPSLocationJSON(json)
	log.Info(json)
	if (valid){

		riderObject.Timestamp = time.Now().UnixNano() // timestamp as soon as we can.
		queue.AddNewPosition(*riderObject)

		// warninglist := queue.RetrieveCollisionList(*riderObject)

		warninglist := queue.RetrieveCollisionList_2(*riderObject,t_param,d_param)
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
	log.Info(json)

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
			UI:      uuid.UUID{},
			Timestamp: time.Now().UnixNano(),
		}

		//timespan := searchObject.Timespan
		//distance := searchObject.Distance

		list := queue.RetrieveCollisionList_2(*GPSSearchObject,searchObject.Timespan,searchObject.Distance,30)
		json := queue.ConvertToMapBoxFreindlyJSON(list)

		if (json == ""){
			log.Info("Invalid json sent from web client")
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w,"server could not parse json parameter")
		}else {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, json)
		}
	}else{
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w,"server not valid json parameter")
	}

}

// to be able to check if it is alive from
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(" HTTP status code returned running: release-1.0.2-demo "))
}


func ServeMap() error {
	http.Handle("/",http.FileServer(http.Dir("./static")));
	return http.ListenAndServe(":8081", nil)
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


    http.HandleFunc("/addposnoret",addNoJsonResponse)
	http.HandleFunc("/addposition", handleNewEntries)
    http.HandleFunc("/retrieve",handleGPSFence)
    http.HandleFunc("/version",pingHandler)
    go log.Fatal(ServeMap())
	log.Fatal(http.ListenAndServe(":8081", nil))

    log.Info("closing server down ... ")
}
