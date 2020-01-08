package main

import (
	"fmt"
	"github.com/google/uuid"
	"os/exec"
	"queue"
	"runtime"
)

const (
	json           = "{}"
	simplecommand  = "http://localhost:8081/version"
	complexcommand = "http://localhost:8081/addposition?gps="
)

func getClimatePayload() *queue.Climatepayload {
	cl := &queue.Climatepayload{
		Ambientemp: 23.3,
		Cabintemp:  19.7,
		Drivertemp: 22.0,
	}

	return cl
}

func getParam1() string {

	payloadstr := queue.GetClimatepayloadJSON(*getClimatePayload())
	gps := &queue.GPSLocation{
		Location: queue.Locationdata{
			Latitude:  37.387401,
			Longitude: -122.035179,
			Accuracy:  1,
			Payload:   payloadstr,
		},
		Gpsobject: queue.Car,
		UI:        uuid.New(),
		Timestamp: 1,
	}

	ajson := queue.GetGPSLocationJSON(*gps)
	return ajson
}

func execute() {

	cmd := exec.Command("ab", "-l", "-n 1000", "-c 100", complexcommand+getParam1())
	//cmd := exec.Command("ab","-n 100","-c 10",simplecommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
}

func main() {

	if runtime.GOOS == "windows" {
		fmt.Println("Can't Execute this on a windows machine")
	} else {
		execute()
	}

}
