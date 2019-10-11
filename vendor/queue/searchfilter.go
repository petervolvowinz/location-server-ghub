package queue

type Filtervalues struct{
	distancevalue int64
	timespanvalue int64
}

type Filter func(a,b,c interface{}) int

type TimeDistFilter struct{
	distance int64
	time int64
}


func GetFilterValues() *Filtervalues{
	return new(Filtervalues)
}

func SetFilterValue(fv *Filtervalues,distancevalue int64,timespanvalue int64) *Filtervalues{
	fv.distancevalue = distancevalue
	fv.timespanvalue = timespanvalue
	return fv
}


