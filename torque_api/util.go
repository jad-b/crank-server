package torque

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// VariadicJoin lets you skip using []string{...} notation.
func VariadicJoin(delim string, args ...string) string {
	return strings.Join(args, delim)
}

// SlashJoin performs a strings.Join using '/' as a separator.
func SlashJoin(args ...string) string {
	return VariadicJoin("/", args...)
}

// PrettyJSON pretty-prints JSON. If an error occurs, you'll get back an empty,
// but valid, JSON structure.
func PrettyJSON(v interface{}) (string, error) {
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return "", err
	}
	return string(s), err
}

//PfmtJSON always returns pretty json or an error.
func PfmtJSON(v interface{}) (s string) {
	var err error
	if s, err = PrettyJSON(v); err != nil {
		return err.Error()
	}
	return s
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
