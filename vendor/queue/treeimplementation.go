package queue

import (
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	log "github.com/sirupsen/logrus"
)

func (t TreeExtended) GetNodeFromKey(key interface{}) (foundNode *rbt.Node){
	node :=  t.tree.Root

	for node != nil {
		compare := t.tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

func IsMemberOf(list [] interface{}, Key GPSLocation)(bool){

	found := false
	i := 0

	for found == false && i < len(list){
		loc := list[i].(GPSLocation)
		found = (loc.Location.Zindex == Key.Location.Zindex) && (loc.UI == Key.UI)
		i++
	}
	return found
}

func (T TreeExtended) AddRoadUserPosition(object interface {}){

	gps := object.(GPSLocation)
	gps.Location.Zindex = GetZorderIndex(gps.Location.Latitude,gps.Location.Longitude)
	T.tree.Put(object,object)
}

// Find predecessor and successor to a tree node. O(h) where the h is height of the tree. h = log n at worst case
func (T *TreeExtended) FindPreSuc(root *rbt.Node,key interface{},pre *rbt.Node,suc *rbt.Node){
	if (root != nil) {

		if T.tree.Comparator(root.Key, key) == 0 {

			// max value in left subtree is predecessor

			if root.Left != nil {
				tmp := root.Left
				for tmp.Right != nil {
					tmp = tmp.Right
				}
				*pre = *tmp

			}

			if root.Right != nil {
				tmp := root.Right
				for tmp.Left != nil {
					tmp = tmp.Left
				}
				*suc = *tmp
			}
			//return pre,suc
		}else if T.tree.Comparator(root.Key, key) == 1 {
			*suc = *root
			T.FindPreSuc(root.Left, key, pre, suc)
		} else {
			*pre = *root
			T.FindPreSuc(root.Right, key, pre, suc)
		}
	}
}

func (T *TreeExtended)  GetNearbyRoadUsers(comparee interface{}, filterdata interface{}, comparator Filter, depth ...int) []interface{}{
	var listofdectees []interface{}

	found := false

	stack := lls.New()
	stack.Push(comparee)

	pre := rbt.Node{}
	var suc rbt.Node

	for found == false  && (stack.Empty() ==  false){

		key,_ := stack.Pop()
		T.FindPreSuc(T.tree.Root,key,&pre,&suc )
		if pre.Key == nil && suc.Key == nil {
			break;
		}

		if pre.Key != nil {
			compresult := comparator(pre.Key, comparee, filterdata)
			if (compresult == 1 ) {
				if (!IsMemberOf(listofdectees, pre.Key.(GPSLocation))) {
					listofdectees = append(listofdectees, pre.Key.(GPSLocation))
					stack.Push(pre.Key)
				}
			}
		}
		if suc.Key != nil {
			compresult := comparator(pre.Key, comparee, filterdata)
			if (compresult == 1) {
				if (!IsMemberOf(listofdectees, suc.Key.(GPSLocation))) {
					listofdectees = append(listofdectees, suc.Key.(GPSLocation))
					stack.Push(suc.Key)
				}
			}
		}


		found = (len(listofdectees) >= depth[0])
	}

	return listofdectees
}

func (T *TreeExtended) GarbageCollect() {
	// start with first
	it := T.tree.Iterator()
	if it.First() {
		firstval := it.Key()
		if (toberetired(firstval.(GPSLocation))){
			T.tree.Remove(firstval)
		}
	}

	// then try last
	it = T.tree.Iterator()
	if it.Last(){
		lastval := it.Key()
		if (toberetired(lastval.(GPSLocation))){
			T.tree.Remove(lastval)
		}
	}

	log.Info("Q SIZE IS : ", T.tree.Size())
}

func (T *TreeExtended) getNearbyRoadUserCandidate(driver GPSLocation,detect GPSLocation,vehicletype int) bool{
	// if it is the same don't add

	if (driver.UI == detect.UI) {
		return false
	}

	// don't need to check type its going to be mutual exclusive
	if (withinTimeSpan(driver.Timestamp,detect.Timestamp,Timedepth)){
		return true
	}

	return false
}

