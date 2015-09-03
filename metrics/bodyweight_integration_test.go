// +build integration
// +build test,db

package integration_test

import (
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/metrics"
)

var (
	bw = metrics.Bodyweight{
		Timestamp: time.Now(),
		Weight:    180.76,
		Comment:   "This is only a test",
	}
	db *sqlx.DB
)

func init() {
	db = torque.OpenDBConnection(torque.LoadPostgresConfig())
	if db == nil {
		log.Fatal("Failed to open database connection. Aborting tests.")
	}
}

func TestDBCreate(t *testing.T) {
	err := bw.Create(db)
	if err != nil {
		t.Error(err)
	}
}

func TestDBRetrieve(t *testing.T) {
	err := bw.Create(db)
	if err != nil {
		t.Error("Failed to create record")
	}
	newBW := metrics.Bodyweight{
		Timestamp: bw.Timestamp,
	}
	err = newBW.Retrieve(db)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(bw, newBW) {
		t.Errorf("Failed to retrieve identical record; %+v vs. %+v", bw, newBW)
	}
}

func TestDBUpdate(t *testing.T) {
	err := bw.Create(db)
	if err != nil {
		t.Error("Failed to create record")
	}
	newBW := metrics.Bodyweight{
		Timestamp: bw.Timestamp,
		Weight:    190.6,
		Comment:   "Nah",
	}
	err = newBW.Update(db)
	if err != nil {
		t.Error(err)
	}
}

func TestDBDelete(t *testing.T) {
	err := bw.Create(db)
	if err != nil {
		t.Error("Failed to create record")
	}
	err = bw.Delete(db)
	if err != nil {
		t.Error(err)
	}
}
