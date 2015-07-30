package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

/*
Package cli is the command-line interface for Torque.


Command syntax:
	torque <options> <action> <resource> arguments>
*/

var (
	host = flag.String("h", "127.0.0.1:8000", "Host:port of Torque server")
)

// Custom timestamp flag
type timestamp time.Time

func (ts *timestamp) String() string {
	return time.Time(*ts).String()
}

func (ts *timestamp) Set(value string) error {
	MyTimeFormat := "2006Jan21504"
	t, err := time.Parse(MyTimeFormat, value)
	if err != nil {
		return err
	}
	*ts = timestamp(t)
	return nil
}

// Bodyweight is a thing
type Bodyweight struct {
	Timestamp time.Time
	Weight    float64
	Comment   string
}

func bwParse(args []string) Bodyweight {
	// Define sub-flags for the bodyweight resource
	var tsFlag timestamp
	bwFlags := flag.NewFlagSet("bwFlags", flag.ContinueOnError)
	bwFlags.Var(&tsFlag, "timestamp", "")
	weight := bwFlags.Float64("weight", 0.0, "")
	comment := bwFlags.String("comment", "", "")

	// Parse the given flags
	bwFlags.Parse(args)

	// Assign their values to a struct
	return Bodyweight{time.Time(tsFlag), *weight, *comment}
}

func main() {
	flag.Parse()
	log.SetOutput(os.Stderr)

	// Handle all errors generically
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s is an invalid call of torque", os.Args)
		}
	}()

	// Check we received a minimal amount of arguments
	remainder := flag.Args()
	lenRemainder := len(remainder)
	log.Printf("Remaining args: %s", remainder)
	if lenRemainder < 1 {
		log.Printf("No action specified")
	} else if lenRemainder < 2 {
		log.Printf("No action specified")
	} else if lenRemainder < 3 {
		log.Printf("No data was provided")
	}
	resource, action := remainder[0], remainder[1]
	switch resource {
	case "bodyweight":
		ret := bwParse(remainder[2:])
		fmt.Printf("%s %+v\n", action, ret)
	default:
		fmt.Printf("%s not recognized as resource", remainder[1])
	}
}
