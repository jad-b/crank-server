package main

import (
	"encoding/json"
	"log"
	"time"
	"testing"
)

func TestTimeFormatting(t *testing.T) {
	t1 := time.Now()
	log.Printf("Time is now %v", t1)
	log.Printf("In RFC3339 format:\n\t%s", t1.Format(time.RFC3339))
	log.Printf("In RFC3339Nano format:\n\t%s", t1.Format(time.RFC3339Nano))
	type TimeStruct struct {
		Timestamp time.Time
	}
	ts1 := &TimeStruct{Timestamp: t1}
	timeAsJSON, _ := json.Marshal(ts1)
	log.Printf("As marshalled JSON:\n\t%s", string(timeAsJSON))
}
