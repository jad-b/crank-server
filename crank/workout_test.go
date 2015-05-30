package crank

import (
	"testing"
	"time"
)

func TestCreateWorkout(t *testing.T) {
	w := NewWorkout()
	if &w.Timestamp == nil {
		t.Error("Timestamp should've defaulted to now.")
	} else if diff := time.Now().Sub(w.Timestamp).Nanoseconds(); diff > 1000 {
		t.Errorf("Timestamp is more than one millisecond old: %d ns", diff)
	}
}
