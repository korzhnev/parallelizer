package parallelizer

type object struct {
	id    int // object ID
	index int // serial number (from 0 to N). Needed for the Union-Find data structure
}

func newObject(id int, index int) object {
	return object{
		id:    id,
		index: index,
	}
}
