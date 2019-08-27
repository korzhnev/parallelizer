package parallelizer

import (
	"testing"
)

type testAction struct {
	id        int
	objectIds []int
	result    chan<- int
}

func (a *testAction) Id() int {
	return a.id
}

func (a *testAction) ObjectIds() []int {
	return a.objectIds
}

func (a *testAction) Work() {
	a.result <- a.id
}

func TestParallelizer_HandleActionsBatch(t *testing.T) {
	actionResult := make(chan int, 4)
	testActions := []testAction{
		{10, []int{200, 300}, actionResult}, // 0
		{11, []int{600, 700}, actionResult}, // 1
		{12, []int{200, 250}, actionResult}, // 2 (linked with 0 by "200")
		{13, []int{700, 750}, actionResult}, // 3 (linked with 1 by "700")
	}
	actions := make([]Action, len(testActions))
	for i, _ := range testActions {
		actions[i] = &testActions[i]
	}

	Run(actions)
	close(actionResult)

	resultOrder := map[int]int{}
	i := 0
	for id := range actionResult {
		resultOrder[id] = i
		i++
	}
	rightOrder := resultOrder[12] > resultOrder[10] && resultOrder[13] > resultOrder[11]
	if !rightOrder {
		t.Error("The order of performing related operations has been violated.")
	}
}
