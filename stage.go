package eskomlol

import "fmt"

type Stage int

// Name returns the descriptive name of the loadshedding stage.
func (s Stage) Name() string {
	switch s {
	case -1:
		return "Unknown"
	case 0:
		return "No Loadshedding"
	default:
		return fmt.Sprintf("Stage %d", s)
	}
}

// Valid determines if the stage is within the valid range.
func (s Stage) Valid() bool {
	var exists bool
	for _, stage := range stageMap {
		if s == stage {
			exists = true
			break
		}
	}
	return exists
}

var stageMap map[int]Stage = map[int]Stage{
	-1: -1,
	1:  0,
	2:  1,
	3:  2,
	4:  3,
	5:  4,
	6:  5,
	7:  6,
	8:  7,
	9:  8,
}
