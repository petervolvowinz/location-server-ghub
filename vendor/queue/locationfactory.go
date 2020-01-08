package queue

import (
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"sync"
)

// var queueMutex = &sync.Mutex{}
var queueOnce sync.Once
var treeOnce sync.Once

/**************************************************************************************
Interface which contains the api to insert,retrieve and remove road user position data.
**************************************************************************************/

type  RoadUsers interface{
	AddRoadUserPosition(val interface{})
	GetNearbyRoadUsers(comparee interface{}, filterdata interface{}, comparator Filter, depth ...int) []interface{}
	GarbageCollect()
	getNearbyRoadUserCandidate()
}

var DataStructureSelection = "T"
var locations RoadUsers

type Queue struct{
	queue *dll.List
}

type TreeExtended struct{
	tree *rbt.Tree
}

func InitRoadUsers() {
	Rf := RoadUsersFactory()
	locations = Rf("")
}

// instance variables
var(
	q_instance  *dll.List
	t_instance *rbt.Tree
)

func GetQueue() *dll.List{
	queueOnce.Do(func(){
		q_instance = dll.New()
	})

	return q_instance
}

func byGPSIndexation(a,b interface {}) int {
	return 0
}

func GetTree() *rbt.Tree{
	treeOnce.Do(func(){
		t_instance = rbt.NewWith(byGPSIndexation)
	})

	return t_instance
}

func RoadUsersFactory() func(datastruct string) RoadUsers{
	return func (datastruct string) RoadUsers {
		var ret RoadUsers
		switch datastruct{
		case "T":
			ret = &TreeExtended{
				tree: GetTree(),
			}
			break
		default:
			ret =  &Queue{
				queue: GetQueue(),
			}
		}
		return ret
	}
}


func withinTimeAndDistanceFilter(a, b, c interface{}) int {

	c1 := a.(GPSLocation)
	c2 := b.(GPSLocation)
	c3 := c.(*TimeDistFilter)

	timespan := c3.time
	distance := c3.distance

	withinTime := withinTimeSpan(c1.Timestamp, c2.Timestamp, timespan)
	if withinTime {
		withinDistance := withinDistance(c1, c2, distance)
		if withinDistance {
			return 1
		} else {
			return 0
		}
	} else {
		return -1
	}

}

func withinTimeSpan(driver_ts int64, detect_ts int64, timespan int64) bool {
	return ((Abs(driver_ts - detect_ts)) / 1e+9) < timespan
}

func withinDistance(driver GPSLocation, detect GPSLocation, distance int64) bool {
	lat1 := driver.Location.Latitude
	long1 := driver.Location.Longitude

	lat2 := detect.Location.Latitude
	long2 := detect.Location.Longitude

	dist := GetApproxDistance2(lat1, long1, lat2, long2)

	return (dist < float64(distance))
}
