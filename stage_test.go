package eskomlol

import "testing"

func TestStage(t *testing.T) {
	s := stageMap[-1]
	if s.Name() != "Unknown" {
		t.Errorf("expected name to be Unknown. got: %s", s.Name())
	}

	s = stageMap[0]
	if s.Name() != "No Loadshedding" {
		t.Errorf("expected name to be No Loadshedding. got: %s", s.Name())
	}

	s = stageMap[8]
	if s.Name() != "Stage 7" {
		t.Errorf("expected name to be Stage 7. got: %s", s.Name())
	}

	fakeStage := Stage(9)
	if fakeStage.Valid() {
		t.Error("expected stage 9 to not be valid")
	}

	fakeStage = Stage(8)
	if !fakeStage.Valid() {
		t.Error("expected stage 8 to be valid")
	}
}
