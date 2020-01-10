package queue



func Add(object interface{}) {
	if locationdata != nil {
		locationdata.AddRoadUserPosition(object)
	}
}

//finds all according to the filter and the comparator function
func FindAll(comparee interface{}, filterdata interface{}, comparator Filter, depth ...int) []interface{} {
	/*var resultarray []interface{}

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

	return resultarray*/
	var breaklimit int = 10 // only return maximum 10 data points for now
	if len(depth) > 0 {
		breaklimit = depth[0]
	}

	return locationdata.GetNearbyRoadUsers(comparee,filterdata,comparator,breaklimit)
}


/*func Find() int {

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
}*/
