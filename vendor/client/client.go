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



func getClimatePayload() *queue.Climatepayload{
	cl := &queue.Climatepayload{
		Ambientemp: 23.3,
		Cabintemp:19.7,
		Drivertemp:22.0,
	}

	return cl;
}

func getSearchParam() string{
	search := &queue.Searchstruct{
		Latitude:37.387401,
		Longitude:-122.035179,
		Distance:200,
		Timespan:5,
	}

	return queue.GetJSONFromGPSLocationObject(*search)
}

func getParam1() string {

	payloadstr := queue.GetClimateJSON(*getClimatePayload())
	gps := &queue.GPSLocation{
		Location: queue.Locationdata{
			Latitude:  37.387401,
			Longitude: -122.035179,
			Accuracy:  1,
			Payload:payloadstr,
		},
		Gpsobject: queue.Car,
		Uuid:      uuid.New(),
		Timestamp: 1,
	}

	ajson := queue.GetJSONFromGPSLocationObject(*gps)
	return ajson
}


/******* JSON returned from addposition for climate data *********
{"Warnings":
[
{"Location":{"lat":37.387401,"lng":-122.035179,"accuracy":1,"payload":"{\"ambientemp\":23.3,\"cabintemp\":19.7,\"dirvertemp\":22}"},"Gpsobject":1,"Uuid":"9dce1043-6c00-4213-beb6-3204efb880c5","Timestamp":1569452142978191000},
{"Location":{"lat":37.387401,"lng":-122.035179,"accuracy":1,"payload":"{\"ambientemp\":23.3,\"cabintemp\":19.7,\"dirvertemp\":22}"},"Gpsobject":1,"Uuid":"7eca7078-0b87-42c1-b9b4-d5826249d6fb","Timestamp":1569452141973377000},
{"Location":{"lat":37.387401,"lng":-122.035179,"accuracy":1,"payload":"{\"ambientemp\":23.3,\"cabintemp\":19.7,\"dirvertemp\":22}"},"Gpsobject":1,"Uuid":"676bb1dc-249f-4f8b-be2e-e7a06016512f","Timestamp":1569452139967853000},
{"Location":{"lat":37.387401,"lng":-122.035179,"accuracy":1,"payload":"{\"ambientemp\":23.3,\"cabintemp\":19.7,\"dirvertemp\":22}"},"Gpsobject":1,"Uuid":"3bdf8375-0383-48f6-8832-2b8c1355cd4c","Timestamp":1569452138961848000},
{"Location":{"lat":37.387401,"lng":-122.035179,"accuracy":1,"payload":"{\"ambientemp\":23.3,\"cabintemp\":19.7,\"dirvertemp\":22}"},"Gpsobject":1,"Uuid":"a89d25af-4276-42f2-a3cb-57764eb98f37","Timestamp":1569452136956759000}
]
}
**********************************************************************/

/************ ADD POSITION CALL *******
http://<IPAdresss>:<POST>/addposition?gps={"Location":{
		"lat":37.387401,
		"lng":-122.035179,
        "accuracy":1,
         "payload":"{\"ambientemp\":23.3,\"cabintemp\":19.7,\"drivertemp\":22}"},
         "Gpsobject":1,
         "Uuid":"8082dc40-6a9f-4eb4-9108-8ae9fcd28cc4",
         "Timestamp":1}
}
********/
func getParam2() string {
	payloadstr := queue.GetClimateJSON(*getClimatePayload())
	gps := &queue.GPSLocation{
		Location: queue.Locationdata{
			Latitude:  37.387401,
			Longitude: -122.035179,
			Accuracy:  1,
			Payload:payloadstr,
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
		resp1, err := http.Get("http://localhost:8081/addposition?gps="+json)
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
		log.Info("addposition ", json)

		resp2, err := http.Get("http://localhost:8081/addposition?gps="+json)
       // resp3, err := http.Get("http://localhost:8081/retrieve?search="+getSearchParam())

       // bytes,_ := ioutil.ReadAll(resp3.Body)
       // log.Info("retrieve ", string(bytes))
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
