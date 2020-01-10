package queue

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

/****** QUEUE IMPLEMENTATION ***************/

func (Q *Queue) AddRoadUserPosition(object interface {}){

	if Q.queue != nil {
		Q.queue.Prepend(object)
	}
}

func (Q *Queue)  GetNearbyRoadUsers(comparee interface{}, filterdata interface{}, comparator Filter, depth ...int) []interface{}{
	var resultarray []interface{}

	var breaklimit int = 10 // only return maximum 10 data points for now
	if len(depth) > 0 {
		breaklimit = depth[0]
	}

	iterator := Q.queue.Iterator()
	iterationcounter := 0
	for iterator.Next() {
		iterationcounter++
		compresult := comparator(iterator.Value(), comparee, filterdata)
		if compresult == 1 {
			resultarray = append(resultarray, iterator.Value())
		} else if compresult == -1 {
			break
		}
		if iterationcounter >= breaklimit {
			break
		}
	}

	return resultarray
}

func toberetired(oldie GPSLocation) bool {

	currentsecond := time.Now().UnixNano()
	oldiesecond := oldie.Timestamp

	log.Println("time checking ", (currentsecond-oldiesecond)/1e+9)
	return (currentsecond-oldiesecond)/1e+9 > Expirationtime

}


func (Q *Queue) GarbageCollect(){
	qInst := Q.queue
	it := qInst.Iterator()

	if it.Last() && toberetired(it.Value().(GPSLocation)) {
		qInst.Remove(qInst.Size() - 1)
	}
	log.Info("Q SIZE IS : ", qInst.Size())

}

func (Q *Queue) getNearbyRoadUserCandidate(driver GPSLocation,detect GPSLocation,vehicletype int) bool{
	fmt.Println("Implementatiom of getNearbyRoadUserCandidat NOT DONE")
	return false
}

