package main

import (
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"net/http"
	"queue"
	"time"
)


func getParam1() string {
	gps := &queue.GPSLocation{
		Location: queue.Locationdata{
			Latitude:  37.387401,
			Longitude: -122.035179,
			Accuracy:  1,
		},
		Gpsobject: queue.Bike,
		Uuid:      uuid.New(),
	}

	ajson := queue.GetJSONFromGPSLocationObject(*gps)
	return ajson
}

/*
{"Location":
   {"lat":37.387401,"lng":-122.035179,"accuracy":1},
 "Gpsobject":0,
  "Uuid":"c6159901-1812-4c99-a0df-e3140303a4d2",
  "Timestamp":1568994267419327600
}
sdd
*/

func getParam2() string {
	gps := &queue.GPSLocation{
		Location: queue.Locationdata{
			Latitude:  37.387401,
			Longitude: -122.035179,
			Accuracy:  1,
		},
		Gpsobject: queue.Car,
		Uuid:      uuid.New(),
		Timestamp: 1,
	}

	ajson := queue.GetJSONFromGPSLocationObject(*gps)
	return ajson
}



func simpleSimulation(ch chan int){
	for {
		json := getParam1()
		resp1, err := http.Get("http://localhost:6060/addposition?gps="+json)
		if err != nil {
			log.Println(err)
		}else {
			bodybytes , err := ioutil.ReadAll(resp1.Body)
			if (err != nil){
				log.Error(" could  not read http body data ", err)
			}
			log.Info(string(bodybytes))
			resp1.Body.Close()
		}

		time.Sleep(time.Second*2)

		json = getParam2()
		resp2, err := http.Get("http://localhost:6060/addposition?gps="+json)

		if err != nil {
			log.Println(err)
		}else {
			bodybytes , err := ioutil.ReadAll(resp2.Body)
			if (err != nil){
				log.Error(" could  not read http body data ", err)
			}
			log.Info(string(bodybytes))
			resp2.Body.Close()
		}
		time.Sleep(1* time.Second)
	}
	ch <- 1 // stop
}


func main(){

	var ch chan int = make(chan int)
	go simpleSimulation(ch)
	log.Info(" Waiting for subscription to end ...")
	fmt.Println(<-ch);

}
