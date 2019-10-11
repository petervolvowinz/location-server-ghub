package queue

import (
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"log"
)


var(
	queueinstance *dll.List
)

// singleton
func GetQueueInstance() *dll.List{

	once.Do(func(){
		queueinstance = dll.New()
	})

	return queueinstance
}


func Add(object interface{}){
	if (queueinstance != nil){
		queueinstance.Prepend(object)
	}
}

func FindAll(comparee interface{},filterdata interface{},comparator Filter)[] interface{}{
	var resultarray [] interface{}

	iterator := queueinstance.Iterator()
	for iterator.Next() { //TODO we actually could stop when we have passed the timespan...
		compresult := comparator(iterator.Value(),comparee,filterdata)
		if (compresult == 1){
			resultarray = append(resultarray, iterator.Value())
		}else if (compresult == -1){ // do not search passed timespan
			break;
		}
	}

	return resultarray
}

func RemoveAll(fromindex int){

	newinstance := queueinstance.Select(func(index int,value interface{}) bool{
		log.Print(" index " , index , " fromindex " , fromindex)
		return index <= fromindex
	})

    if (!newinstance.Empty()){
		queueinstance = newinstance
	}
}

