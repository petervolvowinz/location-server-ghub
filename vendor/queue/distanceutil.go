package queue

import "math"

const (
	R = 6371000 // earth is not round, its elliptical !!! this value is a mean value
)

//go:inline
func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func GetRadians(degree float64) float64 {
	return degree * math.Pi / 180
}

//returns distance between two geo pos in meter: less accurate
func GetApproxDistance1(lat1 float64, long1 float64, lat2 float64, long2 float64) float64 {

	phi1 := GetRadians(lat1)
	delta1 := GetRadians(long1)

	phi2 := GetRadians(lat2)
	delta2 := GetRadians(long2)

	D := delta1 - delta2
	x := D * math.Cos((phi1+phi2)/2)
	y := (phi1 - phi2)
	d := math.Sqrt(math.Pow(x, 2)+math.Pow(y, 2)) * R

	return d
}

//returns distance between two geo pos in meter: more accurate
func GetApproxDistance2(lat1 float64, long1 float64, lat2 float64, long2 float64) float64 {

	phi1 := GetRadians(lat1)
	delta1 := GetRadians(long1)

	phi2 := GetRadians(lat2)
	delta2 := GetRadians(long2)

	dlon := delta1 - delta2
	dlat := phi1 - phi2

	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(phi1)*math.Cos(phi2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c

	return d

}

func  interleave_uint32_with_zeros(input uint32)(uint64)  {
	// make sure word is 64 bits
	word := uint64(input);

	word = (word ^ (word << 16)) & 0x0000ffff0000ffff // 00000000000000001111111111111111
	word = (word ^ (word << 8 )) & 0x00ff00ff00ff00ff;// 00000000111111110000000011111111
	word = (word ^ (word << 4 )) & 0x0f0f0f0f0f0f0f0f;// 00001111000011110000111100001111
	word = (word ^ (word << 2 )) & 0x3333333333333333;// 00110011001100110011001100110011
	word = (word ^ (word << 1 )) & 0x5555555555555555;// 01010101010101010101010101010101

	return word;
}


func GetZorderIndex(lat float64, long float64)(uint64){

	intrlat := uint32( GetAdjustedLatFloat(lat))  // need to trunc to 32 bits
	intrlong := uint32(GetAdjustedLongFloat(long))

	z_index := interleave_uint32_with_zeros(intrlat) | (interleave_uint32_with_zeros(intrlong) << 1)

	return z_index;
}


// preserving 5 decimals means approx 1 meter precision.
func GetAdjustedLatFloat(f float64)(float64){
	return 	(f + 90.0) * 100000
}

func GetAdjustedLongFloat(f float64)(float64){
	return (f + 180.0) * 100000
}
