package parallelizer

// You should implement this interface in client code
type Action interface {
	Id() int          // action ID
	ObjectIds() []int // list of objects involved in the action
	Work()            // action processing in accordance with client business logic
}
