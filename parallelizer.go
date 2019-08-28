package parallelizer

import (
	"sync"
	"github.com/korzhnev/unionfind"
)

// Process actions bundle
// You can call this function many times in your pipeline
func Run(actions []Action) {
	if len(actions) == 0 {
		return
	}
	newBundle(actions).run()
}

type bundle struct {
	actions               []Action
	uf                    *unionfind.UnionFind
	objectsById           map[int]object
	rootsByObjectIds      map[int]int
	actionsByRoot         map[int][]Action
	actionsWithoutObjects []Action
	serialObjectNumber    int
}

func (b *bundle) run() {
	b.initObjectsById()
	b.unionObjects()
	b.extractRootsAsBucketIds()
	b.putActionsIntoBuckets()
	b.runWorkers()
}

func newBundle(actions []Action) *bundle {
	return &bundle{
		actions:          actions,
		objectsById:      make(map[int]object),
		rootsByObjectIds: make(map[int]int),
		actionsByRoot:    make(map[int][]Action),
	}
}

// Init map of Objects by their ids
func (b *bundle) initObjectsById() {
	for _, a := range b.actions {
		for _, objectId := range a.ObjectIds() {
			if _, exist := b.objectsById[objectId]; !exist {
				b.objectsById[objectId] = newObject(objectId, b.generateObjectSerialNumber())
			}
		}
	}
}

// Union objects inside each action inside Union-Find data structure
func (b *bundle) unionObjects() {
	objectCount := len(b.objectsById)
	var err error
	b.uf, err = unionfind.New(objectCount)
	if err != nil {
		panic(err)
	}

	for _, a := range b.actions {
		var firstObject object
		objectIds := a.ObjectIds()
		if len(objectIds) == 0 {
			continue
		}
		for i, objectId := range objectIds {
			if i == 0 {
				firstObject = b.objectsById[objectId]
			} else {
				err := b.uf.Union(firstObject.index, b.objectsById[objectId].index)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (b *bundle) generateObjectSerialNumber() int {
	value := b.serialObjectNumber
	b.serialObjectNumber++
	return value
}

// Get root for each object via Union-Find data structure.
// Root will be as identifier of action bucket.
func (b *bundle) extractRootsAsBucketIds() {
	for _, object := range b.objectsById {
		root, _ := b.uf.Root(object.index)
		b.rootsByObjectIds[object.id] = root
	}
}

// Select actions that need to be performed sequentially and put them in one bucket
func (b *bundle) putActionsIntoBuckets() {
	for _, a := range b.actions {
		objectIds := a.ObjectIds()
		if len(objectIds) == 0 {
			b.actionsWithoutObjects = append(b.actionsWithoutObjects, a)
			continue
		}
		root := b.rootsByObjectIds[objectIds[0]]
		if _, exist := b.actionsByRoot[root]; !exist {
			b.actionsByRoot[root] = make([]Action, 0)
		}
		b.actionsByRoot[root] = append(b.actionsByRoot[root], a)
	}
}

// Run action handlers. One for each bucket in separated goroutine.
func (b *bundle) runWorkers() {
	var wg sync.WaitGroup

	for _, actions := range b.actionsByRoot {
		wg.Add(1)
		go b.runWorker(&wg, actions)
	}

	if len(b.actionsWithoutObjects) > 0 {
		wg.Add(1)
		go b.runWorker(&wg, b.actionsWithoutObjects)
	}

	wg.Wait()
}

// Handle linked action
func (b *bundle) runWorker(wg *sync.WaitGroup, actions []Action) {
	defer wg.Done()
	for _, a := range actions {
		a.Work()
	}
}
