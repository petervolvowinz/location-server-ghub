package GPSGeneration

import (
	"math"
	"math/rand"
)

type Position struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Fence struct {
	NorthEast Position `json:"position"`
	SouthWest Position `json:"position"`
}


/***************************************************************************
Generate n random GPS positions within the rectangle that is formed by p1,p2.
returns array of Position structs
 ***************************************************************************/

func GenerateNGPSPoints(n int,p1 Position,p2 Position) []Position {

	latSpan := math.Abs(p1.Latitude -  p2.Latitude)
	longSpan:= math.Abs(p1.Longitude - p2.Longitude)
	minLat := math.Min(p1.Latitude,p2.Latitude)
	minLong := math.Min(p1.Longitude,p2.Longitude)

	var positions [] Position

	i := 0
	for i < n {
		newLat := minLat + latSpan * rand.Float64()
		newLong := minLong + longSpan * rand.Float64()
		newp := &Position{
			Latitude:newLat,
			Longitude:newLong,
		}
		positions = append(positions, *newp)
		i++
	}

	return positions
}

func IsGeoPosInsideFence(p Position,p1 Position,p2 Position) bool{

	p1latmin := math.Min(p1.Latitude,p2.Latitude)
	p1latmax := math.Max(p1.Latitude,p2.Latitude)

	p2longmin := math.Min(p1.Longitude,p2.Longitude)
	p2longmax := math.Max(p1.Longitude,p2.Longitude)

	p1prim := &Position{
		Latitude:  p1latmin,
		Longitude: p2longmin,
	}
	p2prim := &Position{
		Latitude:  p1latmax,
		Longitude: p2longmax,
	}

	return (p.Latitude > p1prim.Latitude && p.Latitude < p2prim.Latitude) &&
		(p.Longitude > p1prim.Longitude && p.Longitude < p2prim.Longitude)
}

func (F *Fence) IsGeoPositionInsideFence(position Position) bool{
	return (position.Latitude > F.NorthEast.Latitude  && position.Latitude < F.SouthWest.Latitude) &&
		(position.Longitude > F.NorthEast.Longitude && position.Longitude < F.SouthWest.Longitude)
}

func BuildFence(p1 Position,p2 Position) Fence{
	p1latmin := math.Min(p1.Latitude,p2.Latitude)
	p1latmax := math.Max(p1.Latitude,p2.Latitude)

	p2longmin := math.Min(p1.Longitude,p2.Longitude)
	p2longmax := math.Max(p1.Longitude,p2.Longitude)

	ne := Position{
		Latitude:  p1latmin,
		Longitude: p2longmin,
	}

	sw := Position{
		Latitude:  p1latmax,
		Longitude: p2longmax,
	}

	return Fence{
		NorthEast:ne,
		SouthWest:sw,
	}

}