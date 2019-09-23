package main

import (
	"container/list"
	"fmt"
	"queue"
	"time"
	log "github.com/sirupsen/logrus"
	//"fmt"
	"io"
	"net/http"
)


var detectionQueue *list.List

func handler(w http.ResponseWriter, req *http.Request) {

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
	http.HandleFunc("/addposition", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
