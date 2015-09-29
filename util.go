package torque

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SlashJoin performs a strings.Join using '/' as a separator.
func SlashJoin(args ...string) string {
	return strings.Join(args, "/")
}

// PrettyJSON pretty-prints JSON. If an error occurs, you'll get back an empty,
// but valid, JSON structure.
func PrettyJSON(v interface{}) string {
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Print(err)
		return "{}"
	}
	return string(s)
}

// Ascend searches the current directory and all parent directories looking for
// a file. It returns a File pointer if the filename matches, else an error.
func Ascend(filename string) (f *os.File, err error) {
	// Obtain starting point
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// Setup
	dirpath, path := filepath.Dir(cwd), filepath.Join(cwd, filename)
	for {
		// Terminate condition
		if dirpath == filepath.Dir(path) {
			break
		}

		// Search
		f, err = os.Open(path)
		if err == nil {
			// Confirm it's not a directory
			if info, err := f.Stat(); err == nil {
				if !info.IsDir() {
					return f, nil
				}
			}
		}

		// Increment
		dirpath, path = filepath.Dir(dirpath), filepath.Join(dirpath, filename)
	}
	return nil, fmt.Errorf("Unable to find %s", filename)
}

// Stamp properly formats a time.Time the Torque way.
func Stamp(t time.Time) string {
	return t.Format(ValidTimestamps[0])
}
