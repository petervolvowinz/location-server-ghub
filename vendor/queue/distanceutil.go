package queue

import "math"

const(
	R =  6371000 // earth is not round, its elliptical !!! this value is a mean value
)

//go:inline
func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func GetRadians(degree float64) float64{
	return degree*math.Pi/180
}

//returns distance between two geo pos in meter: less accurate
func GetApproxDistance1(lat1 float64,long1 float64,lat2 float64,long2 float64) float64{

    phi1 := GetRadians(lat1)
	delta1 := GetRadians(long1)

	phi2 := GetRadians(lat2)
	delta2 := GetRadians(long2)

	D := delta1 - delta2
	x := D*math.Cos( (phi1 + phi2) / 2)
	y := (phi1 - phi2)
	d := math.Sqrt(math.Pow(x,2) + math.Pow(y,2)) * R

	return d
}

//returns distance between two geo pos in meter: more accurate
func GetApproxDistance2(lat1 float64,long1 float64,lat2 float64,long2 float64) float64 {

	phi1 := GetRadians(lat1)
	delta1 := GetRadians(long1)

	phi2 := GetRadians(lat2)
	delta2 := GetRadians(long2)

	dlon := delta1 - delta2
	dlat := phi1 - phi2

	a :=  math.Pow( math.Sin(dlat / 2),2)  +  math.Cos(phi1) * math.Cos(phi2) * math.Pow(math.Sin(dlon / 2),2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c

	return d

}