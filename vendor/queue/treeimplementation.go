package queue

import "fmt"

func (T *TreeExtended) AddRoadUserPosition(interface {}){
	fmt.Println("Tree implementation of AddRoadUserPosition")
}

func (T *TreeExtended)  GetNearbyRoadUsers(comparee interface{}, filterdata interface{}, comparator Filter, depth ...int) []interface{}{
	fmt.Println("Tree implementation of GetNearbyRoadUsers")
	return nil
}

func (T *TreeExtended) GarbageCollect(){
	fmt.Println("Tree implementation of garbage collect")
}

func (T *TreeExtended) getNearbyRoadUserCandidate(){
	fmt.Println("Tree implementatiom of getNearbyRoadUserCandidat")
}

