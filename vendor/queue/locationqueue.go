package queue

import (
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"log"
)

var (
	queueinstance *dll.List
)

// singleton
func GetQueueInstance() *dll.List {

	once.Do(func() {
		queueinstance = dll.New()
	})

	return queueinstance
}

func Add(object interface{}) {
	if queueinstance != nil {
		queueinstance.Prepend(object)
	}
}

//finds all according to the filter and the comparator function
func FindAll(comparee interface{}, filterdata interface{}, comparator Filter, depth ...int) []interface{} {
	var resultarray []interface{}

	var breaklimit int = 10 // only return maximum 10 data points for now
	if len(depth) > 0 {
		breaklimit = depth[0]
	}

	iterator := queueinstance.Iterator()
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


func Find() int {

	index, _ := queueinstance.Find(func(index int, value interface{}) bool {
		//log.Info(" index ", index)
		val := retieree(value.(GPSLocation))
		return val
	})

	return index
}


func RemoveAll(fromindex int) {

	newinstance := queueinstance.Select(func(index int, value interface{}) bool {
		log.Print(" index ", index, " fromindex ", fromindex)
		return index <= fromindex
	})

	if !newinstance.Empty() {
		queueinstance = newinstance
	}
}
